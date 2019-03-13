-- +migrate Up
ALTER TABLE operations
  ADD COLUMN slug VARCHAR(255) AFTER id
;

UPDATE operations SET slug = id;

ALTER TABLE operations
  MODIFY COLUMN slug VARCHAR(255) NOT NULL,
  ADD UNIQUE INDEX (slug)
;

-- +migrate Down
ALTER TABLE operations
  DROP COLUMN slug
;
