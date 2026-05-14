-- 202512062050_core.up.sql

-- Migration: Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    owner_contact VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

-- Migration: Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

-- Migration: Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    owner_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Migration: Create environments table
CREATE TABLE IF NOT EXISTS environments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,  -- dev-asia, staging, prod-eu
    tier VARCHAR(255) NOT NULL,  -- development, staging, production
    cluster VARCHAR(255) NOT NULL,
    namespace VARCHAR(255) NOT NULL,
    argocd_instance VARCHAR(255) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (project_id, name),
    CHECK (tier IN ('development', 'staging', 'production'))
);

-- Migration: Create services table
CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    repo_url TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (project_id, name),
    UNIQUE (repo_url)
);

CREATE INDEX IF NOT EXISTS idx_users_team_id
    ON users(team_id);

CREATE INDEX IF NOT EXISTS idx_projects_owner_team_id
    ON projects(owner_team_id);

CREATE INDEX IF NOT EXISTS idx_projects_created_by
    ON projects(created_by);

CREATE INDEX IF NOT EXISTS idx_projects_created_at
    ON projects(created_at);

CREATE INDEX IF NOT EXISTS idx_environments_project_tier
    ON environments(project_id, tier);

CREATE INDEX IF NOT EXISTS idx_services_project_created
    ON services(project_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_services_created_by
    ON services(created_by);
