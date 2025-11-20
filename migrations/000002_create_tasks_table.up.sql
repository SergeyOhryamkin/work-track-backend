-- Create track_items table
CREATE TABLE IF NOT EXISTS track_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    emergency_call BOOLEAN NOT NULL DEFAULT FALSE,
    holiday_call BOOLEAN NOT NULL DEFAULT FALSE,
    working_hours DECIMAL(10, 2) NOT NULL DEFAULT 0,
    working_shifts DECIMAL(10, 2) NOT NULL DEFAULT 0,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_track_items_user_id ON track_items(user_id);

-- Create index on date for date range queries
CREATE INDEX IF NOT EXISTS idx_track_items_date ON track_items(date);

-- Create composite index on user_id and date for efficient filtering
CREATE INDEX IF NOT EXISTS idx_track_items_user_date ON track_items(user_id, date);
