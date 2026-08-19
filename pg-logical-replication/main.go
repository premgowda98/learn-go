package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
)

// ChangeEvent represents a single database change
type ChangeEvent struct {
	Operation string                 `json:"operation"` // INSERT, UPDATE, DELETE
	Table     string                 `json:"table"`
	Schema    string                 `json:"schema"`
	Data      map[string]interface{} `json:"data"`
	OldData   map[string]interface{} `json:"old_data"` // For UPDATE/DELETE
	Timestamp time.Time              `json:"timestamp"`
	LSN       string                 `json:"lsn"`
	XID       uint32                 `json:"xid"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Shutting down...")
		cancel()
	}()

	_ = godotenv.Load()

	pgDSN := os.Getenv("PG_CONNECTION_STRING")
	if pgDSN == "" {
		pgDSN = "postgres://postgres:postgres@localhost:5436/test_db?sslmode=disable"
	}
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:19092"
	}
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "pg-replication-events"
	}

	slotName := "my_replication_slot"
	publicationName := "my_publication"

	log.Printf("========================================")
	log.Printf("PostgreSQL CDC (Change Data Capture)")
	log.Printf("========================================")
	log.Printf("Mode: Logical Replication (pgoutput)")
	log.Printf("PostgreSQL: %s", pgDSN)
	log.Printf("Kafka Broker: %s", kafkaBroker)
	log.Printf("Kafka Topic: %s", kafkaTopic)
	log.Printf("Replication Slot: %s", slotName)
	log.Printf("Publication: %s", publicationName)
	log.Printf("========================================")

	// Connect using pgconn directly (required for replication mode)
	connConfig, err := pgconn.ParseConfig(pgDSN)
	if err != nil {
		log.Fatalf("Failed to parse connection config: %v", err)
	}
	connConfig.RuntimeParams["replication"] = "database"

	conn, err := pgconn.ConnectConfig(ctx, connConfig)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close(ctx)
	log.Println("✓ Connected to PostgreSQL (replication mode)")

	// Connect to Kafka
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()
	log.Println("✓ Connected to Kafka")

	if err := runCDC(ctx, conn, kafkaWriter, slotName, publicationName); err != nil {
		if ctx.Err() != nil {
			log.Println("Stopped.")
		} else {
			log.Fatalf("CDC error: %v", err)
		}
	}
}

func runCDC(ctx context.Context, conn *pgconn.PgConn, kafkaWriter *kafka.Writer, slotName, publicationName string) error {
	// Start logical replication using pglogrepl (handles the protocol correctly)
	err := pglogrepl.StartReplication(ctx, conn, slotName, pglogrepl.LSN(0), pglogrepl.StartReplicationOptions{
		PluginArgs: []string{
			"proto_version '1'",
			fmt.Sprintf("publication_names '%s'", publicationName),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to start replication: %w", err)
	}

	log.Printf("✓ Replication started from slot: %s", slotName)
	log.Println("📨 Listening for changes... (Ctrl+C to stop)")

	// Table metadata cache: OID → column info
	relations := make(map[uint32]*pglogrepl.RelationMessage)

	// Track current transaction XID
	var currentXID uint32

	// LSN to acknowledge back to PostgreSQL
	var clientXLogPos pglogrepl.LSN

	// Standby heartbeat ticker (PostgreSQL requires regular keepalives)
	standbyTicker := time.NewTicker(10 * time.Second)
	defer standbyTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-standbyTicker.C:
			// Send keepalive so PostgreSQL knows we're alive and doesn't drop the slot
			err := pglogrepl.SendStandbyStatusUpdate(ctx, conn, pglogrepl.StandbyStatusUpdate{
				WALWritePosition: clientXLogPos,
				ReplyRequested:   false,
			})
			if err != nil {
				return fmt.Errorf("failed to send keepalive: %w", err)
			}
		default:
			// Set a short deadline so we can check the ticker without blocking forever
			receiveCtx, receiveCancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
			rawMsg, err := conn.ReceiveMessage(receiveCtx)
			receiveCancel()

			if err != nil {
				if pgconn.Timeout(err) {
					continue // Timeout is fine, loop back to check ticker
				}
				if ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("receive error: %w", err)
			}

			// pglogrepl messages arrive as CopyData
			copyData, ok := rawMsg.(*pgproto3.CopyData)
			if !ok {
				continue
			}

			switch copyData.Data[0] {
			case pglogrepl.PrimaryKeepaliveMessageByteID:
				// PostgreSQL is asking us to respond
				pka, err := pglogrepl.ParsePrimaryKeepaliveMessage(copyData.Data[1:])
				if err != nil {
					log.Printf("⚠️ Keepalive parse error: %v", err)
					continue
				}
				if pka.ReplyRequested {
					if err := pglogrepl.SendStandbyStatusUpdate(ctx, conn, pglogrepl.StandbyStatusUpdate{
						WALWritePosition: clientXLogPos,
						ReplyRequested:   false,
					}); err != nil {
						return fmt.Errorf("failed to send keepalive reply: %w", err)
					}
				}

			case pglogrepl.XLogDataByteID:
				xld, err := pglogrepl.ParseXLogData(copyData.Data[1:])
				if err != nil {
					log.Printf("⚠️ XLogData parse error: %v", err)
					continue
				}

				// Parse the pgoutput logical replication message
				logicalMsg, err := pglogrepl.Parse(xld.WALData)
				if err != nil {
					log.Printf("⚠️ WAL parse error: %v", err)
					continue
				}

				event, err := handleMessage(logicalMsg, relations, &currentXID, xld.WALStart)
				if err != nil {
					log.Printf("⚠️ Handle error: %v", err)
					continue
				}

				if event != nil {
					eventJSON, _ := json.MarshalIndent(event, "", "  ")
					log.Printf("\n🔔 Change Event:\n%s", eventJSON)

					if err := publishToKafka(ctx, kafkaWriter, event); err != nil {
						log.Printf("⚠️ Kafka publish error: %v", err)
					}
				}

				// Advance our confirmed position
				if xld.WALStart > clientXLogPos {
					clientXLogPos = xld.WALStart
				}
			}
		}
	}
}

// handleMessage converts a pglogrepl message into a ChangeEvent
func handleMessage(msg pglogrepl.Message, relations map[uint32]*pglogrepl.RelationMessage, currentXID *uint32, lsn pglogrepl.LSN) (*ChangeEvent, error) {
	switch m := msg.(type) {
	case *pglogrepl.BeginMessage:
		*currentXID = m.Xid
		log.Printf("  → BEGIN xid=%d", m.Xid)
		return nil, nil

	case *pglogrepl.CommitMessage:
		log.Printf("  → COMMIT xid=%d lsn=%s", *currentXID, lsn)
		return nil, nil

	case *pglogrepl.RelationMessage:
		// Cache table schema so we can resolve column names in Insert/Update/Delete
		relations[m.RelationID] = m
		log.Printf("  → RELATION %s.%s (OID=%d, %d cols)", m.Namespace, m.RelationName, m.RelationID, len(m.Columns))
		return nil, nil

	case *pglogrepl.InsertMessage:
		rel, ok := relations[m.RelationID]
		if !ok {
			return nil, fmt.Errorf("unknown relation OID %d for INSERT", m.RelationID)
		}
		data, err := tupleToMap(m.Tuple, rel)
		if err != nil {
			return nil, err
		}
		return &ChangeEvent{
			Operation: "INSERT",
			Schema:    rel.Namespace,
			Table:     rel.RelationName,
			Data:      data,
			Timestamp: time.Now(),
			LSN:       lsn.String(),
			XID:       *currentXID,
		}, nil

	case *pglogrepl.UpdateMessage:
		rel, ok := relations[m.RelationID]
		if !ok {
			return nil, fmt.Errorf("unknown relation OID %d for UPDATE", m.RelationID)
		}
		newData, err := tupleToMap(m.NewTuple, rel)
		if err != nil {
			return nil, err
		}
		oldData, err := tupleToMap(m.OldTuple, rel)
		if err != nil {
			return nil, err
		}
		event := &ChangeEvent{
			Operation: "UPDATE",
			Schema:    rel.Namespace,
			Table:     rel.RelationName,
			Data:      newData,
			OldData:   oldData,
			Timestamp: time.Now(),
			LSN:       lsn.String(),
			XID:       *currentXID,
		}
		if m.OldTuple != nil {
			oldData, err := tupleToMap(m.OldTuple, rel)
			if err == nil {
				event.OldData = oldData
			}
		}
		return event, nil

	case *pglogrepl.DeleteMessage:
		rel, ok := relations[m.RelationID]
		if !ok {
			return nil, fmt.Errorf("unknown relation OID %d for DELETE", m.RelationID)
		}
		oldData, err := tupleToMap(m.OldTuple, rel)
		if err != nil {
			return nil, err
		}
		return &ChangeEvent{
			Operation: "DELETE",
			Schema:    rel.Namespace,
			Table:     rel.RelationName,
			Data:      oldData,
			Timestamp: time.Now(),
			LSN:       lsn.String(),
			XID:       *currentXID,
		}, nil

	case *pglogrepl.TruncateMessage:
		log.Printf("  → TRUNCATE (%d tables)", len(m.RelationIDs))
		return nil, nil

	case *pglogrepl.TypeMessage:
		return nil, nil

	case *pglogrepl.OriginMessage:
		return nil, nil

	default:
		log.Printf("  → Unknown message type: %T", msg)
		return nil, nil
	}
}

// tupleToMap converts a pglogrepl TupleData into a column name → value map
func tupleToMap(tuple *pglogrepl.TupleData, rel *pglogrepl.RelationMessage) (map[string]interface{}, error) {
	if tuple == nil {
		return nil, nil
	}

	result := make(map[string]interface{}, len(tuple.Columns))
	for i, col := range tuple.Columns {
		if i >= len(rel.Columns) {
			break
		}
		colName := rel.Columns[i].Name

		switch col.DataType {
		case 'n': // NULL
			result[colName] = nil
		case 'u': // Unchanged TOAST
			result[colName] = nil
		case 't': // Text representation
			result[colName] = string(col.Data)
		case 'b': // Binary
			result[colName] = col.Data
		}
	}
	return result, nil
}

// publishToKafka sends a ChangeEvent as JSON to the Kafka topic
func publishToKafka(ctx context.Context, kafkaWriter *kafka.Writer, event *ChangeEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	key := fmt.Sprintf("%s.%s", event.Schema, event.Table)
	err = kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("kafka write failed: %w", err)
	}

	log.Printf("  ✓ Published to Kafka (topic: %s, key: %s)", kafkaWriter.Topic, key)
	return nil
}
