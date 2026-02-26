-- +migrate Up
CREATE TABLE user_operation_permissions (
  user_id INT NOT NULL,
  operation_id INT NOT NULL,
  role VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (user_id, operation_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (operation_id) REFERENCES operations(id)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE user_operation_permissions;
