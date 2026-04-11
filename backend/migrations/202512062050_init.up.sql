-- 202512062050_init.up.sql
-- Migration: Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(32) NOT NULL,  -- platform_admin, org_admin, team_lead, developer, viewer
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

-- Migration: Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    environments TEXT[] NOT NULL,
    status VARCHAR(16) NOT NULL,  -- draft, active, archived, deprecated 
    owner_team VARCHAR(255) NOT NULL,
    repo_url TEXT,
    repo_provider VARCHAR(32),
    owner_contact VARCHAR(255),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

-- Migration: Create plugins table
CREATE TABLE IF NOT EXISTS plugins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    type VARCHAR(16) NOT NULL,
    entrypoint TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    scope VARCHAR(16) NOT NULL, -- global, project, environment
    description TEXT,
    installed_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Migration: Create deployments table
CREATE TABLE IF NOT EXISTS deployments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id),
    environment VARCHAR(32) NOT NULL,
    service VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL,  -- pending, running, completed, failed, rolled_back
    external_ref VARCHAR(255),  -- ArgoCD sync ID
    commit_sha VARCHAR(64),  -- Git commit SHA
    triggered_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    finished_at TIMESTAMP
);

-- Migration: Create scaffold_requests table
CREATE TABLE IF NOT EXISTS scaffold_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id UUID NOT NULL REFERENCES plugins(id),
    template VARCHAR(255) NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    status VARCHAR(16) NOT NULL,  -- pending, approved, running, completed, failed, rejected
    project_id UUID NOT NULL REFERENCES projects(id),
    environment VARCHAR(32) NOT NULL,
    variables JSONB NOT NULL,
    approved_by UUID REFERENCES users(id),
    result_repo_url TEXT,
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Migration: Create refresh_tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
