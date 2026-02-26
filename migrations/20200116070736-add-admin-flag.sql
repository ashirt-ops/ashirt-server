-- +migrate Up
ALTER TABLE `users` ADD COLUMN `admin` BOOLEAN DEFAULT false AFTER `email`;
ALTER TABLE `auth_scheme_data` ADD COLUMN `must_reset_password` BOOLEAN DEFAULT false AFTER `encrypted_password`;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `admin`;
ALTER TABLE `auth_scheme_data` DROP COLUMN `must_reset_password`;
