-- +migrate Up
SET search_path TO authorizer_service;

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL,
    resource TEXT NOT NULL,
    action TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    version INTEGER DEFAULT 1,
    created_by UUID NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (resource, action),

    CONSTRAINT fk_permissions_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE
);

CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_action ON permissions(action);
CREATE INDEX idx_permissions_slug ON permissions(slug);


-- Update updated_at automatically
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_permissions_timestamp
BEFORE UPDATE ON permissions
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();
