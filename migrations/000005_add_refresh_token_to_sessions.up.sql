-- Migration to add refresh_token to user_sessions table
ALTER TABLE user_sessions ADD COLUMN refresh_token TEXT;
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token ON user_sessions(refresh_token);
