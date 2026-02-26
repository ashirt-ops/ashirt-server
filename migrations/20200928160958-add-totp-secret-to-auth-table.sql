-- +migrate Up
ALTER TABLE `auth_scheme_data`
  ADD `totp_secret` VARCHAR(255) AFTER `must_reset_password`
;

-- +migrate Down
ALTER TABLE `auth_scheme_data`
  DROP `totp_secret`
;
