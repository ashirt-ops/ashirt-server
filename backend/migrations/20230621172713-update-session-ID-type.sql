-- +migrate Up
ALTER TABLE sessions
  MODIFY COLUMN id CHAR(43)
;
-- +migrate Down
ALTER TABLE sessions
  MODIFY COLUMN id NOT NULL AUTO_INCREMENT
;
