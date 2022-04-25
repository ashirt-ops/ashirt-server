-- +migrate Up
ALTER TABLE `evidence_metadata`
      ADD COLUMN `work_started_at` TIMESTAMP DEFAULT NULL AFTER `created_at`
    , ADD COLUMN `status` VARCHAR(255) DEFAULT NULL AFTER `body`
    , ADD COLUMN `last_run_message` TEXT DEFAULT NULL AFTER `status`
    , ADD COLUMN `can_process` BOOLEAN DEFAULT NULL AFTER `last_run_message`
;


-- +migrate Down
ALTER TABLE `evidence_metadata`
      DROP COLUMN `work_started_at`
    , DROP COLUMN `status`
    , DROP COLUMN `last_run_message`
    , DROP COLUMN `can_process`
;
