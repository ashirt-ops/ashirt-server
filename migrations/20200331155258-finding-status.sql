-- +migrate Up
ALTER TABLE `findings` ADD COLUMN `ready_to_report` BOOLEAN NOT NULL DEFAULT 0 AFTER `operation_id`;
ALTER TABLE `findings` ADD COLUMN `ticket_link` VARCHAR(255) NULL AFTER `ready_to_report`;

-- +migrate Down
ALTER TABLE `findings` DROP COLUMN `ready_to_report`;
ALTER TABLE `findings` DROP COLUMN `ticket_link`;
