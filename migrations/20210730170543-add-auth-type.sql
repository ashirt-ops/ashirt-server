-- +migrate Up
ALTER TABLE `auth_scheme_data` ADD COLUMN `auth_type` varchar(255) NOT NULL AFTER `auth_scheme`;
UPDATE `auth_scheme_data` SET `auth_type` = `auth_scheme` WHERE `id` > 0;

-- +migrate Down
ALTER TABLE `auth_scheme_data` DROP COLUMN `auth_type`;
