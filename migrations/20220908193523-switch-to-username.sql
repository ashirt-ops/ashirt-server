-- +migrate Up
ALTER TABLE `auth_scheme_data` CHANGE `user_key` `username` VARCHAR(255) NOT NULL;

-- +migrate Down
ALTER TABLE
    `auth_scheme_data` CHANGE `username` `user_key` VARCHAR(255) NOT NULL;
