-- +migrate Up
CREATE TABLE queries (
    `id` INT NOT NULL AUTO_INCREMENT,
    `operation_id` INT NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `query` VARCHAR(255) NOT NULL,
    `type` VARCHAR(15) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE (`name`, `operation_id`, `type`),
    UNIQUE (`query`, `operation_id`, `type`),
    INDEX(`operation_id`)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE queries;
