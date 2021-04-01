-- +migrate Up

CREATE TABLE `finding_categories` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `category` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `category` (`category`)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `finding_categories`
    (`category`)
VALUES
    ('Product'),
    ('Network'),
    ('Enterprise'),
    ('Vendor'),
    ('Behavioral'),
    ('Detection Gap')
;

-- +migrate Down

DROP TABLE `finding_categories`;
