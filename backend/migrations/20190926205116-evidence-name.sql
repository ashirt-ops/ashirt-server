-- +migrate Up
ALTER TABLE events
 MODIFY column `name` TEXT NOT NULL;

-- +migrate Down
ALTER TABLE events
 MODIFY column `name` VARCHAR(255) NOT NULL;
