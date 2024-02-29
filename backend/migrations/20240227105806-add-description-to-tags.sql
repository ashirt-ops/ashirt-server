
-- +migrate Up
ALTER TABLE tags
	ADD COLUMN `description` VARCHAR(150);
-- +migrate Down
ALTER TABLE tags
	DROP COLUMN `description`;
