-- ============================================================
-- PostgreSQL CDC Setup
-- Run this once before starting the Go CDC application.
-- ============================================================


-- Step 1: Verify logical replication is enabled
-- Must return "logical". If not, see POSTGRES_SETUP.md.
SHOW wal_level;


-- Step 2: Create tables
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


-- Step 3: Create publication
-- Defines which tables the CDC app will receive changes for.
CREATE PUBLICATION my_publication FOR TABLE users, orders, products;


-- Step 4: Create replication slot
-- Tracks the consumer's WAL position. PostgreSQL holds WAL until the
-- consumer confirms each message, so no changes are lost on restart.
SELECT pg_create_logical_replication_slot('my_replication_slot', 'pgoutput');


-- ============================================================
-- Verification
-- ============================================================

-- Check wal_level
SHOW wal_level;

-- Check publication
SELECT pubname, pubinsert, pubupdate, pubdelete
FROM pg_publication
WHERE pubname = 'my_publication';

-- Check which tables are in the publication
SELECT schemaname, tablename
FROM pg_publication_tables
WHERE pubname = 'my_publication';

-- Check replication slot
SELECT slot_name, slot_type, active, confirmed_flush_lsn
FROM pg_replication_slots
WHERE slot_name = 'my_replication_slot';


-- ============================================================
-- Replica Identity
-- Controls how much of the old row is sent in UPDATE/DELETE messages.
--
-- DEFAULT (built-in default): only primary key columns in old_data
-- FULL: all columns in old_data — needed for full before/after state in CDC
-- NOTHING: old_data is nil
--
-- Check the current setting for your tables:
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

-- To capture full before-state on UPDATE and DELETE, run:
ALTER TABLE users    REPLICA IDENTITY FULL;
ALTER TABLE orders   REPLICA IDENTITY FULL;
ALTER TABLE products REPLICA IDENTITY FULL;


-- ============================================================
-- Optional: Sample data
-- Note: inserts made before the CDC app connects will NOT appear
-- in the replication stream. Use these to seed initial state only.
-- ============================================================

INSERT INTO users (name, email, age) VALUES
    ('Alice Johnson', 'alice@example.com', 28),
    ('Bob Smith', 'bob@example.com', 35),
    ('Charlie Brown', 'charlie@example.com', 42);

INSERT INTO products (name, price, stock) VALUES
    ('Laptop', 999.99, 5),
    ('Mouse', 29.99, 50),
    ('Keyboard', 79.99, 30);

INSERT INTO orders (user_id, amount, status) VALUES
    (1, 99.99, 'completed'),
    (2, 149.50, 'pending'),
    (1, 75.25, 'completed');
