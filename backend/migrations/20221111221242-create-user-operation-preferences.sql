-- +migrate Up
CREATE TABLE `user_operation_preferences` (
    `user_id` int NOT NULL,
    `operation_id` int NOT NULL,
    `is_favorite` boolean NOT NULL DEFAULT false,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`user_id`, `operation_id`),
    KEY `operation_id` (`operation_id`),
    CONSTRAINT `user_operation_preferences_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `user_operation_preferences_ibfk_2` FOREIGN KEY (`operation_id`) REFERENCES `operations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET = utf8
;


INSERT INTO `user_operation_preferences` (`user_id`, `operation_id`, `is_favorite`)
    SELECT `user_id`, `operation_id`, `is_favorite` from `user_operation_permissions`
;

ALTER TABLE
    `user_operation_permissions` DROP COLUMN `is_favorite`;

-- +migrate Down

ALTER TABLE `user_operation_permissions`
    ADD COLUMN `is_favorite` BOOLEAN DEFAULT false AFTER role
;

UPDATE `user_operation_permissions` `perm`
    INNER JOIN `user_operation_preferences` `pref` ON (
        `perm`.`user_id` = `pref`.`user_id`
        AND `perm`.`operation_id` = `pref`.`operation_id`
    )
    SET `perm`.`is_favorite` = `pref`.`is_favorite`
    WHERE true
;

DROP TABLE `user_operation_preferences`
;
