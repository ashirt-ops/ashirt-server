-- +migrate Up
ALTER TABLE users ADD COLUMN `headless` BOOLEAN DEFAULT false AFTER `admin`;
UPDATE `users` SET `headless` = true WHERE `id` NOT IN (
    SELECT `user_id` FROM `auth_scheme_data`
);

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `headless`;
