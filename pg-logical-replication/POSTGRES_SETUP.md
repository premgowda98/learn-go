# PostgreSQL Setup Guide

Step-by-step instructions for configuring PostgreSQL for logical replication. Run these once before starting the CDC app.

---

## Prerequisites

PostgreSQL must be configured with:

```
wal_level = logical
max_wal_senders >= 1
max_replication_slots >= 1
```

The easiest way is Docker with flags set at startup.

---

## Option A — Docker (recommended)

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

Wait a few seconds, then verify it's running:

```bash
docker logs postgres-cdc | tail -5
```

You should see: `database system is ready to accept connections`

---

## Option B — Existing PostgreSQL installation

1. Edit `postgresql.conf`:

   ```bash
   # macOS (Homebrew)
   nano /usr/local/var/postgresql@16/postgresql.conf

   # Linux
   sudo nano /etc/postgresql/16/main/postgresql.conf
   ```

2. Set these values:

   ```
   wal_level = logical
   max_wal_senders = 4
   max_replication_slots = 4
   ```

3. Restart PostgreSQL:

   ```bash
   # macOS
   brew services restart postgresql@16

   # Linux
   sudo systemctl restart postgresql
   ```

---

## Database setup

Run the provided SQL file to create tables, publication, and replication slot:

```bash
psql -h localhost -U postgres -p 5436 -d test_db -f setup.sql
```

Or run the commands manually:

### 1. Verify logical replication is enabled

```sql
SHOW wal_level;
-- Expected: logical
```

If you see `replica` or `minimal`, stop and fix the PostgreSQL config first.

### 2. Create tables

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    age INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. Create a publication

A publication defines which tables to replicate:

```sql
CREATE PUBLICATION my_publication FOR TABLE users, orders, products;
```

Verify:

```sql
SELECT pubname, puballtables, pubinsert, pubupdate, pubdelete
FROM pg_publication
WHERE pubname = 'my_publication';

SELECT schemaname, tablename
FROM pg_publication_tables
WHERE pubname = 'my_publication';
```

### 4. Create a replication slot

The slot tracks how far the consumer has read. PostgreSQL holds WAL until the slot confirms it's been consumed:

```sql
SELECT pg_create_logical_replication_slot('my_replication_slot', 'pgoutput');
```

Verify:

```sql
SELECT slot_name, slot_type, active, restart_lsn, confirmed_flush_lsn
FROM pg_replication_slots
WHERE slot_name = 'my_replication_slot';
```

`active = f` means no consumer is connected yet (expected at this point).

---

## Verify the full setup

```sql
-- wal_level must be "logical"
SHOW wal_level;

-- Publication exists
SELECT pubname FROM pg_publication WHERE pubname = 'my_publication';

-- Tables are in the publication
SELECT tablename FROM pg_publication_tables WHERE pubname = 'my_publication';

-- Replication slot exists
SELECT slot_name, active FROM pg_replication_slots WHERE slot_name = 'my_replication_slot';
```

All four queries should return rows. If any are empty, re-run the corresponding setup step.

---

## Test data

Insert some rows to verify the setup before running the CDC app:

```sql
INSERT INTO users (name, email, age) VALUES
    ('Alice Johnson', 'alice@example.com', 28),
    ('Bob Smith', 'bob@example.com', 35);

INSERT INTO products (name, price, stock) VALUES
    ('Laptop', 999.99, 5),
    ('Mouse', 29.99, 50);

INSERT INTO orders (user_id, amount, status) VALUES
    (1, 99.99, 'completed'),
    (2, 149.50, 'pending');
```

These inserts happened before the CDC app was running, so they won't appear in the stream. Only changes made after the app connects will be captured. (This is expected behavior — the slot tracks the WAL position from when the consumer connects.)

---

## Monitoring while the app runs

Open a second psql session and run these to watch replication progress:

```sql
-- Is the slot active? Is it keeping up?
SELECT
    slot_name,
    active,
    confirmed_flush_lsn,
    pg_current_wal_lsn(),
    (pg_current_wal_lsn() - confirmed_flush_lsn) AS lag_bytes
FROM pg_replication_slots
WHERE slot_name = 'my_replication_slot';
```

- `active = t` — the CDC app is connected
- `lag_bytes = 0` — the app is keeping up
- `lag_bytes` growing — the app is falling behind

```sql
-- Active replication connections
SELECT pid, usename, application_name, state, sent_lsn, write_lsn, flush_lsn
FROM pg_stat_replication;
```

---

## Teardown

To reset and start fresh:

```sql
-- Drop the replication slot (required before dropping tables/publication)
SELECT pg_drop_replication_slot('my_replication_slot');

-- Drop the publication
DROP PUBLICATION IF EXISTS my_publication;

-- Drop the tables
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;
```

---

## Troubleshooting

### `wal_level is not logical`

You need to set `wal_level = logical` in `postgresql.conf` and restart PostgreSQL. If using Docker, stop the container and recreate it with the `-c wal_level=logical` flag.

### `replication slot "my_replication_slot" does not exist`

Run:
```sql
SELECT pg_create_logical_replication_slot('my_replication_slot', 'pgoutput');
```

### `publication "my_publication" does not exist`

Run:
```sql
CREATE PUBLICATION my_publication FOR TABLE users, orders, products;
```

### `replication slot ... is active for PID ...`

Another consumer process is connected. Find and stop it, or create a second slot with a different name.

### App connects but no messages appear

- Insert data in PostgreSQL **after** the app starts (pre-existing data won't appear)
- Check that the table is in the publication: `SELECT * FROM pg_publication_tables WHERE pubname = 'my_publication';`
- Check that `wal_level = logical` is still set (it resets on some managed databases)
