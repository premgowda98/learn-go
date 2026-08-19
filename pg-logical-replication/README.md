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

## Replica Identity and old_data

PostgreSQL controls how much of the old row is included in UPDATE and DELETE messages via a per-table setting called **REPLICA IDENTITY**.

| Mode | What's in `old_data` on UPDATE / DELETE |
|------|------------------------------------------|
| `DEFAULT` (default) | Primary key columns only — all other columns are null |
| `FULL` | Every column with its previous value |
| `NOTHING` | Nothing — `old_data` is nil |
| `USING INDEX` | Columns covered by a specific unique index |

### INSERT vs UPDATE vs DELETE

| Operation | `data` (new row) | `old_data` (previous row) |
|-----------|-----------------|--------------------------|
| INSERT | Full row, always | Not applicable |
| UPDATE | Full new row, always | Depends on REPLICA IDENTITY |
| DELETE | Not applicable | Depends on REPLICA IDENTITY |

This means:
- With `DEFAULT`: UPDATE gives you the full new row but only the primary key in `old_data`. You can't see what the previous values of non-key columns were.
- With `FULL`: UPDATE gives you both the complete before and after state — the standard behaviour expected from CDC tools like Debezium.

### Checking current replica identity

```sql
SELECT
    relname AS table_name,
    CASE relreplident
        WHEN 'd' THEN 'DEFAULT (primary key only)'
        WHEN 'f' THEN 'FULL (all columns)'
        WHEN 'n' THEN 'NOTHING'
        WHEN 'i' THEN 'USING INDEX'
    END AS replica_identity
FROM pg_class
WHERE relname IN ('users', 'orders', 'products')
  AND relkind = 'r';
```

### Enabling full old-row capture

To get the complete previous row in UPDATE and DELETE events, run:

```sql
ALTER TABLE users    REPLICA IDENTITY FULL;
ALTER TABLE orders   REPLICA IDENTITY FULL;
ALTER TABLE products REPLICA IDENTITY FULL;
```

After this, `old_data` in UPDATE events will contain every column's value before the change.

### Seeing exactly what changed — `changed_fields`

On UPDATE events, the app computes a field-level diff between `old_data` and `data` and includes it as `changed_fields`. Each entry shows the before and after value for only the columns that actually changed:

```json
{
  "operation": "UPDATE",
  "table": "users",
  "schema": "public",
  "data": {
    "id": "1",
    "name": "Alice Updated",
    "email": "alice@example.com",
    "age": "31"
  },
  "old_data": {
    "id": "1",
    "name": "Alice",
    "email": "alice@example.com",
    "age": "30"
  },
  "changed_fields": {
    "name": { "from": "Alice", "to": "Alice Updated" },
    "age":  { "from": "30",   "to": "31" }
  },
  "lsn": "0/1C4B2A0",
  "xid": 751
}
```

`changed_fields` is only populated when `REPLICA IDENTITY FULL` is set — without it, `old_data` only has the primary key, so a meaningful diff isn't possible.

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
    Operation     string                 `json:"operation"`                // INSERT | UPDATE | DELETE
    Table         string                 `json:"table"`
    Schema        string                 `json:"schema"`
    Data          map[string]interface{} `json:"data"`                     // new row values
    OldData       map[string]interface{} `json:"old_data,omitempty"`       // previous row (requires REPLICA IDENTITY FULL)
    ChangedFields map[string]FieldDiff   `json:"changed_fields,omitempty"` // diff, UPDATE only
    Timestamp     time.Time              `json:"timestamp"`
    LSN           string                 `json:"lsn"`
    XID           uint32                 `json:"xid"`
}

// FieldDiff is one entry in ChangedFields — the before and after value of a column
type FieldDiff struct {
    From interface{} `json:"from"`
    To   interface{} `json:"to"`
}
```

Column values come as text strings (PostgreSQL's text representation). Parse them as needed downstream.

`ChangedFields` is computed automatically on UPDATE by diffing `OldData` against `Data`. It is nil when `REPLICA IDENTITY FULL` is not set (because `OldData` won't have the full previous row to compare against).

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
