
-- +migrate Up
ALTER TABLE default_tags
	ADD COLUMN `description` VARCHAR(100);
-- +migrate Down
ALTER TABLE default_tags
	DROP COLUMN `description`;
