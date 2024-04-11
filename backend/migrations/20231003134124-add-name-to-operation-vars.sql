-- +migrate Up
ALTER TABLE operation_vars
ADD COLUMN name VARCHAR(255) NOT NULL AFTER slug;
-- +migrate Down
ALTER TABLE operation_vars
DROP COLUMN name;

