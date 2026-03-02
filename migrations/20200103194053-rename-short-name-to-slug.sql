-- +migrate Up
DROP INDEX `short_name_identity_provider` ON `users`;
ALTER TABLE `users` CHANGE `short_name` `slug` VARCHAR(255) NOT NULL;
ALTER TABLE `users` ADD UNIQUE INDEX (`slug`);

-- +migrate Down
ALTER TABLE `users`
  CHANGE `slug` `short_name` VARCHAR(255) NOT NULL
;
