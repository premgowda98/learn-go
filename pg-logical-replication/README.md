# PostgreSQL CDC (Change Data Capture) with Kafka

A production-grade Go application that reads database changes directly from PostgreSQL's logical replication stream and publishes them to Kafka in real-time. No polling. No SELECT queries. Pure streaming CDC.

---

## How It Works

```
PostgreSQL (INSERT / UPDATE / DELETE)
    ↓
Write-Ahead Log (WAL)
    ↓
Replication Slot  ←── pgoutput protocol
    ↓
Go App (pglogrepl)
    ├── Receives binary WAL messages
    ├── Parses pgoutput protocol (Relation, Insert, Update, Delete, ...)
    ├── Resolves column names from cached relation metadata
    └── Publishes JSON ChangeEvents to Kafka
    ↓
Kafka Topic: pg-replication-events
```

PostgreSQL writes every change to its Write-Ahead Log (WAL). A logical replication slot decodes that WAL into structured messages using the `pgoutput` plugin. This app connects to that slot and streams those messages in real-time — exactly how tools like Debezium work.

---

## Quick Start

### 1. Start PostgreSQL with logical replication enabled

```bash
docker run --name postgres-cdc \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=test_db \
  -p 5436:5432 \
  -d postgres:16 \
  -c wal_level=logical \
  -c max_wal_senders=4 \
  -c max_replication_slots=4
```

### 2. Set up the database

```bash
psql -h localhost -U postgres -p 5436 -d test_db -f setup.sql
```

### 3. Start Kafka (Redpanda, single-node)

```bash
docker run --name redpanda \
  -p 19092:19092 \
  -d redpandadata/redpanda:latest \
  redpanda start \
  --kafka-addr PLAINTEXT://0.0.0.0:19092 \
  --advertise-kafka-addr PLAINTEXT://localhost:19092
```

### 4. Configure the app

```bash
cp .env.template .env
# Edit .env if your ports/credentials differ
```

### 5. Run the CDC app

```bash
go run main.go
```

### 6. Insert data and watch it flow

```bash
psql -h localhost -U postgres -p 5436 -d test_db \
  -c "INSERT INTO users (name, email, age) VALUES ('Alice', 'alice@example.com', 30);"
```

The app will print:

```
🔔 Change Event:
{
  "operation": "INSERT",
  "table": "users",
  "schema": "public",
  "data": {
    "age": "30",
    "email": "alice@example.com",
    "id": "1",
    "name": "Alice",
    ...
  },
  "timestamp": "2026-08-19T21:30:00Z",
  "lsn": "0/1C4A1F0",
  "xid": 742
}
  ✓ Published to Kafka (topic: pg-replication-events, key: public.users)
```

---

## pgoutput Protocol

PostgreSQL sends binary messages over the replication stream. Each message has a single-byte type prefix:

| Byte | Type     | Description                                      |
|------|----------|--------------------------------------------------|
| `B`  | Begin    | Start of a transaction                           |
| `R`  | Relation | Table schema (columns, types) — sent once per table per session |
| `I`  | Insert   | New row data                                     |
| `U`  | Update   | Old row (if REPLICA IDENTITY FULL) + new row     |
| `D`  | Delete   | Old row / key columns                            |
| `C`  | Commit   | End of a transaction                             |
| `T`  | Truncate | Table truncated                                  |
| `O`  | Origin   | Replication origin info                          |

A single `INSERT INTO users ...` produces this sequence:

```
B (Begin, xid=742)
R (Relation: public.users, 6 columns)  ← only on first insert in session
I (Insert: new tuple with all column values)
C (Commit)
```

---

## Implementation

### Key packages

| Package | Role |
|---------|------|
| `github.com/jackc/pglogrepl` | Starts replication, receives and parses pgoutput messages |
| `github.com/jackc/pgx/v5/pgconn` | Low-level PostgreSQL connection in replication mode |
| `github.com/jackc/pgx/v5/pgproto3` | Wire protocol types (CopyData, etc.) |
| `github.com/segmentio/kafka-go` | Kafka producer |

### Connection setup

Logical replication requires a special connection parameter. We use `pgconn` directly (not the higher-level `pgx.Conn`) because replication connections don't support the extended query protocol:

```go
connConfig.RuntimeParams["replication"] = "database"
conn, _ := pgconn.ConnectConfig(ctx, connConfig)
```

### Starting replication

```go
pglogrepl.StartReplication(ctx, conn, slotName, pglogrepl.LSN(0), pglogrepl.StartReplicationOptions{
    PluginArgs: []string{
        "proto_version '1'",
        "publication_names 'my_publication'",
    },
})
```

`LSN(0)` means "start from the earliest unconfirmed position in the slot". In production you'd persist the last confirmed LSN and resume from there.

### Receiving messages

```go
rawMsg, _ := conn.ReceiveMessage(ctx)
copyData := rawMsg.(*pgproto3.CopyData)

switch copyData.Data[0] {
case pglogrepl.XLogDataByteID:
    xld, _ := pglogrepl.ParseXLogData(copyData.Data[1:])
    logicalMsg, _ := pglogrepl.Parse(xld.WALData)
    // handle logicalMsg...

case pglogrepl.PrimaryKeepaliveMessageByteID:
    // reply if server requests it
}
```

### Relation cache

PostgreSQL sends a `RelationMessage` the first time a table appears in the stream. We cache it by OID so we can resolve column names in subsequent Insert/Update/Delete messages:

```go
relations := make(map[uint32]*pglogrepl.RelationMessage)

case *pglogrepl.RelationMessage:
    relations[m.RelationID] = m

case *pglogrepl.InsertMessage:
    rel := relations[m.RelationID]
    // use rel.Columns[i].Name to map column values
```

### Standby keepalives

PostgreSQL expects periodic heartbeats. If we go silent too long, the server closes the connection:

```go
pglogrepl.SendStandbyStatusUpdate(ctx, conn, pglogrepl.StandbyStatusUpdate{
    WALWritePosition: clientXLogPos,
})
```

We send one every 10 seconds, and immediately when the server sets `ReplyRequested = true` on a keepalive.

---

## ChangeEvent structure

Every INSERT, UPDATE, or DELETE produces a `ChangeEvent`:

```go
type ChangeEvent struct {
    Operation string                 `json:"operation"` // INSERT | UPDATE | DELETE
    Table     string                 `json:"table"`
    Schema    string                 `json:"schema"`
    Data      map[string]interface{} `json:"data"`      // new row (or deleted row for DELETE)
    OldData   map[string]interface{} `json:"old_data,omitempty"` // previous row for UPDATE
    Timestamp time.Time              `json:"timestamp"`
    LSN       string                 `json:"lsn"`
    XID       uint32                 `json:"xid"`
}
```

Column values come as text strings (PostgreSQL's text representation). Parse them as needed downstream.

---

## Testing

```bash
# Insert
psql -h localhost -U postgres -p 5436 -d test_db \
  -c "INSERT INTO users (name, email, age) VALUES ('Bob', 'bob@example.com', 25);"

# Update
psql -h localhost -U postgres -p 5436 -d test_db \
  -c "UPDATE users SET age = 26 WHERE email = 'bob@example.com';"

# Delete
psql -h localhost -U postgres -p 5436 -d test_db \
  -c "DELETE FROM users WHERE email = 'bob@example.com';"
```

### Consume from Kafka

```bash
# Using kafka-console-consumer (if you have Kafka CLI)
kafka-console-consumer \
  --bootstrap-server localhost:19092 \
  --topic pg-replication-events \
  --from-beginning

# Using rpk (Redpanda CLI)
rpk topic consume pg-replication-events --brokers localhost:19092
```

---

## File structure

```
.
├── main.go              # CDC consumer — replication stream → Kafka
├── go.mod / go.sum      # Dependencies
├── setup.sql            # PostgreSQL: tables, publication, replication slot
├── .env.template        # Environment variable template
├── POSTGRES_SETUP.md    # Detailed PostgreSQL setup walkthrough
└── README.md            # This file
```

---

## Environment variables

```bash
PG_CONNECTION_STRING=postgres://postgres:postgres@localhost:5436/test_db?sslmode=disable
KAFKA_BROKER=localhost:19092
KAFKA_TOPIC=pg-replication-events
```

---

## PostgreSQL prerequisites

| Setting | Required value | Why |
|---------|---------------|-----|
| `wal_level` | `logical` | Enables logical decoding of WAL |
| `max_wal_senders` | ≥ 1 | Allows replication connections |
| `max_replication_slots` | ≥ 1 | Allows creating slots |

Verify:

```sql
SHOW wal_level;           -- must be "logical"
SHOW max_wal_senders;     -- must be > 0
SHOW max_replication_slots; -- must be > 0
```

A replication slot and publication must exist before starting the app:

```sql
CREATE PUBLICATION my_publication FOR TABLE users, orders, products;
SELECT pg_create_logical_replication_slot('my_replication_slot', 'pgoutput');
```

---

## What's not included (next steps for production)

| Feature | Description |
|---------|-------------|
| LSN persistence | Save `clientXLogPos` to disk/DB on each commit so the app resumes from the right position after restart |
| Kafka retries | Retry failed publishes with exponential backoff before exiting |
| Reconnection | Auto-reconnect to PostgreSQL on network failure |
| Monitoring | Track replication lag (`pg_current_wal_lsn() - confirmed_flush_lsn`), message rates, errors |
| Schema changes | Handle `ALTER TABLE` (relation cache invalidation) |
| Type coercion | Parse text column values into proper Go types based on column OID |

---

## References

- [PostgreSQL Logical Replication](https://www.postgresql.org/docs/current/logical-replication.html)
- [pgoutput Protocol Reference](https://www.postgresql.org/docs/current/protocol-logical-replication.html)
- [jackc/pglogrepl](https://github.com/jackc/pglogrepl) — Go library used for replication
- [jackc/pgx](https://github.com/jackc/pgx) — PostgreSQL driver
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go) — Kafka client
- [Debezium](https://debezium.io/) — Production CDC tool (Java), same approach
