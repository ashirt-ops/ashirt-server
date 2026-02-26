-- +migrate Up
ALTER TABLE ashirt_auth CHANGE email username VARCHAR(255) NOT NULL;
ALTER TABLE users ADD email VARCHAR(255) NOT NULL AFTER last_name;

UPDATE users SET email=(SELECT username FROM ashirt_auth WHERE ashirt_auth.user_id=users.id) WHERE identity_provider='ashirt';

-- +migrate Down
ALTER TABLE ashirt_auth CHANGE username email VARCHAR(255) NOT NULL;
ALTER TABLE users DROP COLUMN email;
