-- +migrate Up
CREATE TABLE api_keys (
  id INT AUTO_INCREMENT,
  user_id INT NOT NULL,
  access_key VARCHAR(255) NOT NULL,
  secret_key VARBINARY(255) NOT NULL,
  last_auth TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  UNIQUE (access_key)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE api_keys;
