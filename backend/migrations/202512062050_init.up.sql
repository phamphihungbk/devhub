-- 202512062050_init.up.sql
-- Migration: Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(16) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL,
    last_login TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Migration: Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    environments TEXT[],
    created_by UUID NOT NULL,
    deleted_at TIMESTAMP
);

-- Migration: Create plugins table
CREATE TABLE IF NOT EXISTS plugins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    type VARCHAR(16) NOT NULL,
    description TEXT,
    installed_at TIMESTAMP NOT NULL
);

-- Migration: Create deployments table
CREATE TABLE IF NOT EXISTS deployments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    environment VARCHAR(32) NOT NULL,
    service VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL,
    triggered_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- Migration: Create scaffold_requests table
CREATE TABLE IF NOT EXISTS scaffold_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template VARCHAR(255) NOT NULL,
    project_id UUID NOT NULL,
    environment VARCHAR(32) NOT NULL,
    variables JSONB NOT NULL
);
