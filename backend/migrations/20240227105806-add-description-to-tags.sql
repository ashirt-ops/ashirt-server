
-- +migrate Up
ALTER TABLE tags
	ADD COLUMN `description` VARCHAR(100);
-- +migrate Down
ALTER TABLE tags
	DROP COLUMN `description`;
