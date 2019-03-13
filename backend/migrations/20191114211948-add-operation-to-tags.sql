-- +migrate Up
-- Add column and remove the unique constraint on name
ALTER TABLE tags
  ADD COLUMN operation_id int AFTER id,
  DROP KEY name
;

-- Insert existing tags for every existing operation_id
INSERT INTO tags
  (name, color_name, created_at, updated_at, operation_id)
  SELECT tags.name, tags.color_name, tags.created_at, tags.updated_at, operations.id AS operation_id
  FROM tags LEFT JOIN operations
  ON true
;

-- Replace tags on evidence with new row with current operation_id
UPDATE tag_evidence_map dest
  LEFT JOIN (
    SELECT tag_evidence_map.tag_id AS old_tag_id, tag_evidence_map.evidence_id, newtags.id AS new_tag_id
    FROM tag_evidence_map
    LEFT JOIN tags AS oldtags ON oldtags.id = tag_id
    LEFT JOIN evidence ON evidence.id = evidence_id
    LEFT JOIN tags AS newtags ON oldtags.name = newtags.name
    WHERE newtags.operation_id = evidence.operation_id
  ) source
  ON source.evidence_id = dest.evidence_id AND source.old_tag_id = dest.tag_id
  SET dest.tag_id = source.new_tag_id
;

-- Delete the original tags
DELETE FROM tags WHERE operation_id IS NULL;

-- Add unique name+operation_id constraint
ALTER TABLE tags
  MODIFY COLUMN operation_id int NOT NULL,
  ADD UNIQUE (`name`, `operation_id`),
  ADD INDEX(`operation_id`)
;

-- +migrate Down
-- Remove unique constraint on (name, operation_id)
ALTER TABLE tags DROP KEY `name`;

-- Drop operation id
ALTER TABLE `tags` DROP COLUMN `operation_id`;

-- Replace tags on evidence
UPDATE tag_evidence_map
  LEFT JOIN (
    SELECT MIN(id) AS new_id, GROUP_CONCAT(id) AS old_ids, name
    FROM tags GROUP BY name
  ) tags
  ON FIND_IN_SET(tag_evidence_map.tag_id, tags.old_ids)
  SET tag_id = tags.new_id
;

-- Remove duplicates
DELETE t2 FROM tags t1
  INNER JOIN tags t2
  WHERE t1.id < t2.id AND t1.name = t2.name
;

-- Add unique constraint back on name
ALTER TABLE tags ADD UNIQUE (`name`);
