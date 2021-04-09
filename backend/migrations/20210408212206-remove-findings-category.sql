-- +migrate Up

ALTER TABLE `findings` ADD COLUMN `category_id` INT NULL DEFAULT NULL AFTER `ticket_link`;
ALTER TABLE `findings` ADD CONSTRAINT `fk_category_id__finding_categories_id` FOREIGN KEY (`category_id`) REFERENCES `finding_categories`(id);

INSERT INTO `finding_categories` (`category`) 
    SELECT DISTINCT `category` FROM `findings` WHERE `category` NOT IN (
        SELECT `category` FROM `finding_categories`
    ) AND `category` != ''
;

UPDATE `findings` SET `category_id` = (
    SELECT `id` FROM `finding_categories` WHERE `category` = `findings`.`category`
)
WHERE `findings`.`category` != ''
;

ALTER TABLE `findings` DROP COLUMN `category`;

-- +migrate Down

ALTER TABLE `findings` ADD COLUMN `category` varchar(255) NOT NULL DEFAULT '' AFTER `category_id`;

UPDATE `findings` SET `category` = (
    SELECT `category` FROM `finding_categories` WHERE `id` = `findings`.`category_id`
)
WHERE `findings`.`category_ID` IS NOT NULL
;

ALTER TABLE `findings` DROP CONSTRAINT `fk_category_id__finding_categories_id`;
ALTER TABLE `findings` DROP COLUMN  `category_id`;
