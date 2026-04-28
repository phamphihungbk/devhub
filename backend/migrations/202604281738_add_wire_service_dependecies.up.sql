-- 202604281738_add_wire_service_dependecies.up.sql

CREATE TABLE service_dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    depends_on_service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,

    type VARCHAR(32) NOT NULL, -- http, grpc, queue, database
    protocol VARCHAR(32),      -- http, https, grpc
    port INT,
    path TEXT,

    config JSONB NOT NULL DEFAULT '{}',

    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    CHECK (service_id <> depends_on_service_id),
    UNIQUE(service_id, depends_on_service_id, type)
);

CREATE INDEX idx_service_dependencies_service_id
    ON service_dependencies(service_id);

CREATE INDEX idx_service_dependencies_depends_on_service_id
    ON service_dependencies(depends_on_service_id);
