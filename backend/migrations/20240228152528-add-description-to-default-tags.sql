
-- +migrate Up
ALTER TABLE default_tags
	ADD COLUMN `description` VARCHAR(150);
-- +migrate Down
ALTER TABLE default_tags
	DROP COLUMN `description`;
