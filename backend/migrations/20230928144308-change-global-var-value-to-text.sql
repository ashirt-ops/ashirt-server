-- +migrate Up
ALTER TABLE global_vars
MODIFY COLUMN value TEXT;
-- +migrate Down
ALTER TABLE global_vars
MODIFY COLUMN value VARCHAR(255);
