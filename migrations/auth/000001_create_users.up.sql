-- Auth Service Migration: 000001_create_users
-- Creates the users table with all required fields and indexes.

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'user',
    avatar        VARCHAR(500),
    bio           VARCHAR(1000),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

-- Indexes for frequent lookup patterns
CREATE INDEX IF NOT EXISTS idx_users_email      ON users(email)      WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_username   ON users(username)   WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
