-- +migrate Up
ALTER TABLE api_keys MODIFY access_key VARBINARY(255) NOT NULL;
-- +migrate Down
ALTER TABLE api_keys MODIFY access_key VARCHAR(255) NOT NULL;
