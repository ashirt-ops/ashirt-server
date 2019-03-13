-- +migrate Up
CREATE TABLE `tag_evidence_map` (
  `tag_id` INT NOT NULL,
  `evidence_id` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`tag_id`,`evidence_id`),
  KEY `evidence_id` (`evidence_id`),
  FOREIGN KEY (`tag_id`) REFERENCES tags(`id`),
  FOREIGN KEY (`evidence_id`) REFERENCES evidence(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE `tag_evidence_map`;
