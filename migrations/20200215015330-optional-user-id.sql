-- +migrate Up
ALTER TABLE `sessions`
  MODIFY `user_id` INT
;

-- +migrate Down
ALTER TABLE `sessions`
  MODIFY `user_id` INT NOT NULL
;
