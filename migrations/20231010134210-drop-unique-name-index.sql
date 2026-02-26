-- +migrate Up
ALTER TABLE operation_vars
DROP INDEX name;
-- +migrate Down
ALTER TABLE operation_vars
ADD UNIQUE INDEX name (slug);

