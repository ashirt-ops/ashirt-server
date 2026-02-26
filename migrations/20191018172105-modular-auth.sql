-- +migrate Up
CREATE TABLE `ashirt_auth` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `email` varchar(255) NOT NULL,
    `encrypted_password` varbinary(255) DEFAULT NULL,
    `user_id` int(11) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO ashirt_auth (email, encrypted_password, user_id)
    SELECT  email, encrypted_password, id FROM users;

ALTER TABLE `users` DROP INDEX `email`;
ALTER TABLE `users` DROP INDEX `short_name`;
ALTER TABLE `users` DROP `email`;
ALTER TABLE `users` DROP `encrypted_password`;

ALTER TABLE `users` ADD `identity_provider` VARCHAR(255) NOT NULL AFTER `last_name`;
ALTER TABLE `users` ADD UNIQUE `short_name_identity_provider`(`short_name`, `identity_provider`);

UPDATE users SET identity_provider = 'ashirt_local';

-- +migrate Down
ALTER TABLE `users` DROP INDEX `short_name_identity_provider`;
ALTER TABLE `users` DROP `identity_provider`;

ALTER TABLE `users` ADD `email` varchar(255) NOT NULL AFTER `last_name`;
ALTER TABLE `users` ADD `encrypted_password` varbinary(255) DEFAULT NULL AFTER `email`;

UPDATE `users` SET 
    `email` = (SELECT email FROM ashirt_auth WHERE user_id = users.id),
    `encrypted_password` = (SELECT encrypted_password FROM ashirt_auth WHERE user_id = users.id);

ALTER TABLE `users` ADD UNIQUE `email` (`email`);
ALTER TABLE `users` ADD UNIQUE `short_name` (`short_name`);

DROP TABLE ashirt_auth;
