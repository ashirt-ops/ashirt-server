-- +migrate Up
ALTER TABLE `users` ADD COLUMN `disabled` BOOLEAN DEFAULT false AFTER `admin`;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE IF NOT EXISTS `sessions` (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    session_data LONGBLOB, 
    created_at TIMESTAMP DEFAULT NOW(), 
    modified_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP, 
    expires_at TIMESTAMP DEFAULT NOW(), 
    PRIMARY KEY(`id`),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `disabled`;
DROP TABLE `sessions`;
