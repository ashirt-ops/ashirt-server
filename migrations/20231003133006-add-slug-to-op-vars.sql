-- +migrate Up
ALTER TABLE operation_vars
CHANGE COLUMN name slug VARCHAR(255) NOT NULL UNIQUE;
-- +migrate Down
ALTER TABLE operation_vars
CHANGE COLUMN slug name VARCHAR(255) NOT NULL UNIQUE;
