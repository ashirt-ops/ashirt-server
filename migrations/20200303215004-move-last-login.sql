-- +migrate Up
alter table users drop column `last_login`;
alter table auth_scheme_data add column `last_login` timestamp default null after must_reset_password;

-- +migrate Down
alter table auth_scheme_data drop column `last_login`;
alter table users add column `last_login` timestamp default null after `disabled`;
