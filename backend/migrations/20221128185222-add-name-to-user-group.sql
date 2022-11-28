-- +migrate Up
ALTER TABLE user_groups
ADD slug VARCHAR(255) UNIQUE;

-- +migrate Down
ALTER TABLE user_groups
REMOVE slug;
