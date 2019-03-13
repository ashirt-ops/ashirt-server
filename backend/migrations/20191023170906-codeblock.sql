-- +migrate Up
ALTER TABLE `evidence` ADD `content_type`  varchar(31) NOT NULL AFTER `description`;
UPDATE `evidence` SET `content_type`='image';

-- +migrate Down
ALTER TABLE `evidence` DROP `content_type`;
