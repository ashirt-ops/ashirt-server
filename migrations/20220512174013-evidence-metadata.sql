-- +migrate Up
CREATE TABLE evidence_metadata (
    `id` INT AUTO_INCREMENT,
    `evidence_id` INT NOT NULL,
    `source` VARCHAR(255) NOT NULL,
    `body` TEXT NOT NULL,
    `status` VARCHAR(255) DEFAULT NULL,
    `last_run_message` TEXT DEFAULT NULL,
    `can_process` BOOLEAN DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `work_started_at` TIMESTAMP DEFAULT NULL,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`evidence_id`) REFERENCES evidence(id),
    UNIQUE(`evidence_id`, `source`),
    INDEX(`evidence_id`)
) ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;

-- +migrate Down
DROP TABLE evidence_metadata;
