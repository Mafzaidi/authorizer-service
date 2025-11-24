-- +migrate Down
SET search_path TO authorizer_service;

DROP TRIGGER IF EXISTS update_users_timestamp ON users;
DROP FUNCTION IF EXISTS update_timestamp();
DROP TABLE IF EXISTS users;
