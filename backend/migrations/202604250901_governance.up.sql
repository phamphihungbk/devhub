-- 202604250901_governance.up.sql

-- Migration: Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Migration: Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Migration: Create role_permissions table
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (role_id, permission_id)
);

-- Migration: Create user_roles table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (user_id, role_id)
);

-- Migration: Create approval_policies table
CREATE TABLE IF NOT EXISTS approval_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(64) NOT NULL,  -- deployment, scaffold_request, release (future)
    service_id UUID REFERENCES services(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    required_approvals INT NOT NULL DEFAULT 1,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (required_approvals > 0)
);

-- Migration: Create approval_requests table
CREATE TABLE IF NOT EXISTS approval_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(64) NOT NULL,
    resource_id UUID NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    status VARCHAR(32) NOT NULL,
    required_approvals INT NOT NULL,
    approved_count INT NOT NULL,
    rejected_count INT NOT NULL,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (status IN ('pending', 'approved', 'rejected')),
    CHECK (required_approvals > 0),
    CHECK (approved_count >= 0),
    CHECK (rejected_count >= 0)
);

-- Migration: Create approval_decisions table
CREATE TABLE IF NOT EXISTS approval_decisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    approval_request_id UUID NOT NULL REFERENCES approval_requests(id) ON DELETE CASCADE,
    decided_by UUID NOT NULL REFERENCES users(id),
    decision VARCHAR(16) NOT NULL,  -- approve, reject
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (decision IN ('approve', 'reject')),
    UNIQUE (approval_request_id, decided_by)
);


CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id
    ON role_permissions(permission_id);

CREATE INDEX IF NOT EXISTS idx_user_roles_role_id
    ON user_roles(role_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_approval_policies_unique_scope
    ON approval_policies(resource, service_id, environment_id)
    NULLS NOT DISTINCT;

CREATE INDEX IF NOT EXISTS idx_approval_policies_lookup
    ON approval_policies(
        resource,
        environment_id,
        enabled,
        service_id,
        updated_at DESC,
        created_at DESC
    );

CREATE INDEX IF NOT EXISTS idx_approval_requests_status_created
    ON approval_requests(status, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_approval_requests_requested_by_created
    ON approval_requests(requested_by, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_approval_requests_resource
    ON approval_requests(resource, resource_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_approval_requests_pending_resource
    ON approval_requests(resource, resource_id)
    WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_approval_decisions_request_created
    ON approval_decisions(approval_request_id, created_at ASC);

CREATE INDEX IF NOT EXISTS idx_approval_decisions_decided_by_created
    ON approval_decisions(decided_by, created_at DESC);
