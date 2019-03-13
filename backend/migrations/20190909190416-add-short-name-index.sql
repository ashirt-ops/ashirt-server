-- +migrate Up
ALTER TABLE users
  MODIFY COLUMN short_name VARCHAR(255) NOT NULL,
  ADD UNIQUE INDEX (short_name)
;

-- +migrate Down
