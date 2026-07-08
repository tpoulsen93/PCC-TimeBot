-- Additive migration: extend the SMS schema to support the web app.
-- Safe to run against a live SMS database — existing data is untouched.
-- Run via: psql "$DATABASE_URL" -f scripts/migrate-add-web-auth.sql

-- Add is_admin to employees (SMS app ignores it; web app uses it for admin routes).
ALTER TABLE employees
    ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE;

-- Magic-link login tokens (web app only).
CREATE TABLE IF NOT EXISTS login_tokens (
    id          SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS login_tokens_employee_idx ON login_tokens (employee_id);

-- Web sessions (web app only).
CREATE TABLE IF NOT EXISTS sessions (
    id          TEXT PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS sessions_employee_idx ON sessions (employee_id);
