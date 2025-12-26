-- Create user_sessions table to track login/logout activity
CREATE TABLE IF NOT EXISTS user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    
    -- Session timestamps
    login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    logout_at TIMESTAMP,
    
    -- Device/platform information
    device TEXT,           -- e.g., "iPhone 14 Pro", "Chrome on Windows", "Android App"
    platform TEXT,         -- e.g., "ios", "android", "web"
    user_agent TEXT,       -- Full user agent string for detailed analysis
    ip_address VARCHAR(45), -- IPv4 or IPv6
    
    -- Session metadata
    session_duration_seconds INTEGER, -- Calculated on logout
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_login_at ON user_sessions(login_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_platform ON user_sessions(platform);
