-- +migrate Up
ALTER TABLE `users` ADD COLUMN `deleted_at` TIMESTAMP DEFAULT NULL AFTER `updated_at`;
ALTER TABLE `auth_scheme_data` ADD CONSTRAINT `fk_user_id__users_id` FOREIGN KEY (user_id) REFERENCES users(id);

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `deleted_at`;
ALTER TABLE `auth_scheme_data` DROP FOREIGN KEY `fk_user_id__users_id`;
