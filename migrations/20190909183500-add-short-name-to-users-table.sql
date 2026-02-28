-- +migrate Up
ALTER TABLE users
  ADD COLUMN short_name VARCHAR(255) AFTER id
;

-- +migrate Down
ALTER TABLE users
  DROP COLUMN short_name
;
