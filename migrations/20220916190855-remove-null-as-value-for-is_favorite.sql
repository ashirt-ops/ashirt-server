-- +migrate Up
ALTER TABLE user_operation_permissions CHANGE
  is_favorite
  is_favorite BOOLEAN NOT NULL DEFAULT '0';

-- +migrate Down
ALTER TABLE user_operation_permissions CHANGE
  is_favorite
  is_favorite BOOLEAN NULL DEFAULT '0';
