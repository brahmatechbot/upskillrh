CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email CITEXT NOT NULL,
  full_name TEXT NOT NULL,
  status TEXT NOT NULL,
  last_login_at TIMESTAMPTZ NULL,
  session_version INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ NULL,
  CONSTRAINT users_status_check CHECK (status IN ('pending', 'active', 'suspended', 'disabled'))
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_active_uq ON users (email) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS user_password_credentials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL UNIQUE REFERENCES users(id),
  password_hash TEXT NOT NULL,
  algorithm TEXT NOT NULL DEFAULT 'argon2id',
  password_changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  requires_change BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT password_algorithm_check CHECK (algorithm IN ('argon2id'))
);

CREATE TABLE IF NOT EXISTS auth_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  token_hash BYTEA NOT NULL UNIQUE,
  csrf_hash BYTEA NOT NULL,
  session_version INTEGER NOT NULL,
  remember_me BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  idle_expires_at TIMESTAMPTZ NOT NULL,
  absolute_expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS auth_sessions_user_active_idx ON auth_sessions (user_id, absolute_expires_at) WHERE revoked_at IS NULL;
