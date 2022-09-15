-- +migrate Up
ALTER TABLE `auth_scheme_data` ADD COLUMN `json_data` json AFTER `totp_secret`;


-- +migrate Down
ALTER TABLE `auth_scheme_data` DROP COLUMN `json_data`;
