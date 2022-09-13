-- +migrate Up
ALTER TABLE user_operation_permissions
  ADD COLUMN is_favorite BOOLEAN DEFAULT false AFTER role
;

-- +migrate Down
ALTER TABLE user_operation_permissions
  DROP COLUMN is_favorite
;
