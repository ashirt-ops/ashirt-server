-- +migrate Up
ALTER TABLE events
  ADD COLUMN uuid VARCHAR(36) AFTER id
;

ALTER TABLE evidence
  ADD COLUMN uuid VARCHAR(36) AFTER id
;

UPDATE events SET uuid = uuid();
UPDATE evidence SET uuid = uuid();

ALTER TABLE events
  MODIFY COLUMN uuid VARCHAR(36) NOT NULL,
  ADD UNIQUE INDEX (uuid)
;

ALTER TABLE evidence
  MODIFY COLUMN uuid VARCHAR(36) NOT NULL,
  ADD UNIQUE INDEX (uuid)
;

-- +migrate Down
ALTER TABLE events
  DROP COLUMN uuid
;

ALTER TABLE evidence
  DROP COLUMN uuid
;
