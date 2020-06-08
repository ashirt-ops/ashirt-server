-- +migrate Up
CREATE TABLE exports_queue (
  `id` INT AUTO_INCREMENT,
  `operation_id` INT NOT NULL,
  `user_id` INT NOT NULL,
  `export_name` VARCHAR(255),
  `status` INT NOT NULL,
  `notes` TEXT,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY `exports_queue_operation_id__operations_id` (`operation_id`) REFERENCES operations(`id`),
  FOREIGN KEY `exports_queue_user_id__users_id` (`user_id`) REFERENCES users(`id`)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE exports_queue;
