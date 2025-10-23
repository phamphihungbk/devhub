-- Migration: Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    environments TEXT[],
    created_by UUID NOT NULL,
    deleted_at TIMESTAMP
);
