-- +migrate Up
CREATE TABLE user_group_operation_permissions (
  group_id INT NOT NULL,
  operation_id INT NOT NULL,
  role VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (group_id, operation_id),
  FOREIGN KEY (group_id) REFERENCES user_groups(id),
  FOREIGN KEY (operation_id) REFERENCES operations(id)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE user_group_operation_permissions;
