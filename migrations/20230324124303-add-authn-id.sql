-- +migrate Up
ALTER TABLE auth_scheme_data
  ADD COLUMN authn_id VARCHAR(255) AFTER user_id
;

-- +migrate Down
ALTER TABLE auth_scheme_data
  DROP COLUMN authn_id
;
