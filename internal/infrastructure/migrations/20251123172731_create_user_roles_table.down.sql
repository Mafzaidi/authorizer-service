-- +migrate Down
SET search_path TO authorizer_service;

DROP TABLE IF EXISTS user_roles;
