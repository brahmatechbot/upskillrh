DROP INDEX IF EXISTS auth_sessions_user_active_idx;
DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS user_password_credentials;
DROP INDEX IF EXISTS users_email_active_uq;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS citext;
