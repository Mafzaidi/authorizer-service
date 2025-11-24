-- +migrate Down
SET search_path TO authorizer_service;

DROP TRIGGER IF EXISTS update_applications_timestamp ON applications;
DROP FUNCTION IF EXISTS update_timestamp();
DROP TABLE IF EXISTS applications;
