-- 202604201016_add_rbac_service.up.sql

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS approval_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(64) NOT NULL,
    action VARCHAR(64) NOT NULL,
    project_id UUID NULL REFERENCES projects(id) ON DELETE CASCADE,
    service_id UUID NULL REFERENCES services(id) ON DELETE CASCADE,
    environment VARCHAR(64) NULL,
    required_approvals INT NOT NULL DEFAULT 1,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (required_approvals > 0)
);


CREATE TABLE IF NOT EXISTS approval_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(64) NOT NULL,
    action VARCHAR(64) NOT NULL,
    resource_id UUID NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    project_id UUID NULL REFERENCES projects(id),
    service_id UUID NULL REFERENCES services(id),
    environment VARCHAR(64) NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    required_approvals INT NOT NULL DEFAULT 1,
    approved_count INT NOT NULL DEFAULT 0,
    rejected_count INT NOT NULL DEFAULT 0,
    resolved_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (status IN ('pending', 'approved', 'rejected', 'canceled')),
    CHECK (required_approvals > 0),
    CHECK (approved_count >= 0),
    CHECK (rejected_count >= 0)
);


CREATE TABLE IF NOT EXISTS approval_decisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    approval_request_id UUID NOT NULL REFERENCES approval_requests(id) ON DELETE CASCADE,
    decided_by UUID NOT NULL REFERENCES users(id),
    decision VARCHAR(16) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (decision IN ('approve', 'reject')),
    UNIQUE(approval_request_id, decided_by)
);

CREATE INDEX IF NOT EXISTS idx_approval_policies_resource_action
    ON approval_policies(resource, action);

CREATE INDEX IF NOT EXISTS idx_approval_policies_scope
    ON approval_policies(project_id, service_id, environment);

CREATE INDEX IF NOT EXISTS idx_approval_requests_requested_by
    ON approval_requests(requested_by);

CREATE INDEX IF NOT EXISTS idx_approval_requests_status
    ON approval_requests(status);

CREATE INDEX IF NOT EXISTS idx_approval_requests_scope
    ON approval_requests(project_id, service_id, environment);

CREATE INDEX IF NOT EXISTS idx_approval_requests_resource
    ON approval_requests(resource, action, resource_id);

CREATE INDEX IF NOT EXISTS idx_approval_decisions_request
    ON approval_decisions(approval_request_id);
