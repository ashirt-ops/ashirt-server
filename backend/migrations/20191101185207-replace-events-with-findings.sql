-- +migrate Up
RENAME TABLE events TO findings;

ALTER TABLE findings
  CHANGE COLUMN `name` `title` VARCHAR(255) NOT NULL,
  ADD COLUMN `description` TEXT NOT NULL AFTER title,
  ADD COLUMN `category` VARCHAR(255) NOT NULL DEFAULT "" AFTER operation_id,

  DROP FOREIGN KEY `findings_ibfk_2`,
  DROP COLUMN `occurred_at`,
  DROP COLUMN `operator_id`
;

DROP TABLE tag_event_map;

UPDATE queries SET type = 'findings' WHERE type = 'events';

RENAME TABLE evidence_event_map TO evidence_finding_map;

ALTER TABLE evidence_finding_map
  CHANGE COLUMN `event_id` `finding_id` INT NOT NULL
;

-- +migrate Down
RENAME TABLE findings TO events;

ALTER TABLE events
  ADD COLUMN `operator_id` int(11) NOT NULL AFTER operation_id,
  ADD COLUMN `occurred_at` timestamp NOT NULL DEFAULT now() AFTER description
;

ALTER TABLE events
  CHANGE COLUMN `title` `name` TEXT,
  DROP COLUMN `description`,
  DROP COLUMN `category`
;

-- (From schema.sql)
CREATE TABLE `tag_event_map` (
  `tag_id` int(11) NOT NULL,
  `event_id` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`tag_id`,`event_id`),
  KEY `event_id` (`event_id`),
  CONSTRAINT `tag_event_map_ibfk_1` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`),
  CONSTRAINT `tag_event_map_ibfk_2` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

UPDATE queries SET type = 'events' WHERE type = 'findings';

RENAME TABLE evidence_finding_map TO evidence_event_map;

ALTER TABLE evidence_event_map
  CHANGE COLUMN `finding_id` `event_id` INT NOT NULL
;
