-- PCC-TimeBot schema (fresh install)
-- Run via: psql "$DATABASE_URL" -f scripts/init-db.sql

DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS login_tokens CASCADE;
DROP TABLE IF EXISTS payroll CASCADE;
DROP TABLE IF EXISTS employees CASCADE;

CREATE TABLE employees (
    id            SERIAL PRIMARY KEY,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    phone         TEXT,
    email         TEXT,
    supervisor_id INTEGER REFERENCES employees(id),
    is_admin      BOOLEAN NOT NULL DEFAULT FALSE,
    timestamp     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX employees_email_unique
    ON employees (lower(email))
    WHERE email IS NOT NULL;

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

CREATE TABLE login_tokens (
    id          SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX login_tokens_employee_idx ON login_tokens (employee_id);

CREATE TABLE sessions (
    id          TEXT PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX sessions_employee_idx ON sessions (employee_id);

INSERT INTO employees (first_name, last_name, email, is_admin)
VALUES ('taylor', 'poulsen', 't-poulsen@hotmail.com', TRUE);
