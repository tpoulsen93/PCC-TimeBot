-- PCC-TimeBot schema (SMS/Twilio version)
-- Run via: psql "$DATABASE_URL" -f scripts/init-db.sql

DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS login_tokens CASCADE;
DROP TABLE IF EXISTS payroll CASCADE;
DROP TABLE IF EXISTS employees CASCADE;

-- Employees
CREATE TABLE employees (
    id            SERIAL PRIMARY KEY,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    phone         TEXT,
    email         TEXT,
    supervisor_id INTEGER REFERENCES employees(id),
    timestamp     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payroll entries
-- The UNIQUE constraint named "submission" is relied on by ON CONFLICT ON CONSTRAINT submission
CREATE TABLE payroll (
    transaction_id SERIAL PRIMARY KEY,
    id             INTEGER NOT NULL REFERENCES employees(id),
    time           DOUBLE PRECISION NOT NULL,
    date           DATE NOT NULL,
    message        TEXT,
    location       TEXT,
    timestamp      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT submission UNIQUE (id, date)
);
