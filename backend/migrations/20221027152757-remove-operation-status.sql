-- +migrate Up
ALTER TABLE operations
  DROP COLUMN status;

-- +migrate Down
ALTER TABLE operations
  ADD COLUMN status INT NOT NULL;
