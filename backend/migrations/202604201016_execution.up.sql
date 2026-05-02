-- 202604201016_execution.up.sql

-- Migration: Create releases table
CREATE TABLE IF NOT EXISTS releases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    tag VARCHAR(64) NOT NULL,  -- v1.2.3
    target VARCHAR(255) NOT NULL,  -- branch or commit
    name VARCHAR(255) NOT NULL,  -- release title
    status VARCHAR(32) NOT NULL,  -- pending, running, completed, failed
    notes TEXT NOT NULL DEFAULT '',
    html_url TEXT NOT NULL,  -- link to GitHub/GitLab release
    external_ref VARCHAR(255) NOT NULL,  -- release ID in SCM
    triggered_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (service_id, tag),
    CHECK (status IN ('pending', 'running', 'completed', 'failed'))
);

-- Migration: Create plugins table
CREATE TABLE IF NOT EXISTS plugins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    type VARCHAR(32) NOT NULL,  -- scaffold, deployment, release
    runtime VARCHAR(32) NOT NULL,  -- node, go, python
    entrypoint TEXT NOT NULL,  -- command / script path
    config_schema JSONB,  -- input validation schema
    description TEXT,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (name, version)
);

-- Migration: Create jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(64) NOT NULL,  -- scaffold, deployment, release, sync_env
    status VARCHAR(32) NOT NULL DEFAULT 'queued',
    resource_type VARCHAR(64) NOT NULL,  -- scaffold_request, deployment, release
    resource_id UUID NOT NULL,
    plugin_id UUID NOT NULL REFERENCES plugins(id),
    payload JSONB NOT NULL DEFAULT '{}',
    result JSONB NOT NULL DEFAULT '{}',
    error TEXT,
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    UNIQUE (type, resource_type, resource_id),
    CHECK (status IN ('queued', 'completed', 'running', 'failed'))
);

-- Migration: Create deployments table
CREATE TABLE IF NOT EXISTS deployments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    release_id UUID NOT NULL REFERENCES releases(id) ON DELETE RESTRICT,
    status VARCHAR(32) NOT NULL,  -- pending, running, completed, failed
    external_ref VARCHAR(255),  -- ArgoCD app / sync ID
    commit_sha VARCHAR(64),  -- Git commit
    triggered_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    CHECK (status IN ('pending', 'running', 'completed', 'failed'))
);

-- Migration: Create scaffold_requests table
CREATE TABLE IF NOT EXISTS scaffold_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    requested_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    approved_by UUID REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(16) NOT NULL,  -- pending, approved, running, completed, failed, rejected
    variables JSONB NOT NULL DEFAULT '{}',  -- user input
    result_repo_url TEXT,  -- created repo
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (status IN ('pending', 'approved', 'running', 'completed', 'failed', 'rejected'))
);

-- Migration: Create refresh_tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    UNIQUE (token)
);

CREATE INDEX IF NOT EXISTS idx_plugins_created_at
    ON plugins(created_at);

CREATE INDEX IF NOT EXISTS idx_jobs_queued_type_created
    ON jobs(type, status, created_at, id)
    WHERE status = 'queued';

CREATE INDEX IF NOT EXISTS idx_jobs_resource
    ON jobs(resource_type, resource_id);

CREATE INDEX IF NOT EXISTS idx_jobs_plugin_id
    ON jobs(plugin_id);

CREATE INDEX IF NOT EXISTS idx_jobs_created_by
    ON jobs(created_by);

CREATE INDEX IF NOT EXISTS idx_deployments_service_created
    ON deployments(service_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_deployments_environment_created
    ON deployments(environment_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_deployments_triggered_by
    ON deployments(triggered_by);

CREATE INDEX IF NOT EXISTS idx_deployments_status_created
    ON deployments(status, created_at, id);

CREATE INDEX IF NOT EXISTS idx_scaffold_requests_project_created
    ON scaffold_requests(project_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_scaffold_requests_requested_by
    ON scaffold_requests(requested_by);

CREATE INDEX IF NOT EXISTS idx_scaffold_requests_status_created
    ON scaffold_requests(status, created_at, id);

CREATE INDEX IF NOT EXISTS idx_releases_service_created
    ON releases(service_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_releases_triggered_by
    ON releases(triggered_by);

CREATE INDEX IF NOT EXISTS idx_releases_status_created
    ON releases(status, created_at, id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_active_expiry
    ON refresh_tokens(user_id, deleted_at, expires_at);
