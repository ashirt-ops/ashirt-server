-- +migrate Up
ALTER TABLE `ashirt_auth`
  CHANGE `username` `user_key` VARCHAR(255),
  MODIFY `user_id` int(11) AFTER `user_key`,
  ADD COLUMN `auth_scheme` VARCHAR(255) AFTER `id`,
  DROP KEY `email`,
  ADD UNIQUE KEY `auth_scheme_user_key` (`auth_scheme`, `user_key`)
;

UPDATE `ashirt_auth` SET `auth_scheme` = "local";

ALTER TABLE `ashirt_auth`
  MODIFY `auth_scheme` VARCHAR(255) NOT NULL
;

RENAME TABLE `ashirt_auth` TO `auth_scheme_data`;

ALTER TABLE `users` DROP COLUMN `identity_provider`;

-- +migrate Down
RENAME TABLE `auth_scheme_data` TO `ashirt_auth`;

ALTER TABLE `ashirt_auth`
  CHANGE `user_key` `username` VARCHAR(255),
  MODIFY `user_id` int(11) AFTER `encrypted_password`,
  DROP COLUMN `auth_scheme`,
  DROP KEY `auth_scheme_user_key`,
  ADD UNIQUE KEY `email` (`username`)
;

ALTER TABLE `users` ADD COLUMN `identity_provider` VARCHAR(255);
UPDATE `users` SET `identity_provider` = "ashirt";
ALTER TABLE `users` MODIFY `identity_provider` VARCHAR(255) NOT NULL AFTER `email`;
