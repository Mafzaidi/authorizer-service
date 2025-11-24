-- +migrate Down
SET search_path TO authorizer_service;

DROP TRIGGER IF EXISTS update_roles_timestamp ON roles;
DROP TABLE IF EXISTS roles;