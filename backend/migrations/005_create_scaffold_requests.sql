-- Migration: Create scaffold_requests table
CREATE TABLE IF NOT EXISTS scaffold_requests (
    id SERIAL PRIMARY KEY,
    template VARCHAR(255) NOT NULL,
    project_id UUID NOT NULL,
    environment VARCHAR(32) NOT NULL,
    variables JSONB NOT NULL
);
