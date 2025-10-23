-- Migration: Create plugins table
CREATE TABLE IF NOT EXISTS plugins (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    type VARCHAR(16) NOT NULL,
    description TEXT,
    installed_at TIMESTAMP NOT NULL
);
