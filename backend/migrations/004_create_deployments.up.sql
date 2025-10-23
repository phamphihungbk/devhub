-- Migration: Create deployments table
CREATE TABLE IF NOT EXISTS deployments (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    environment VARCHAR(32) NOT NULL,
    service VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL,
    triggered_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
);
