-- +migrate Up
CREATE TABLE operation_vars (
  id INT AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,
  value VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


CREATE TABLE var_operation_map (
  operation_id INT NOT NULL,
  var_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (operation_id, var_id),
  FOREIGN KEY (operation_id) REFERENCES operations(id),
  FOREIGN KEY (var_id) REFERENCES operation_vars(id) ON DELETE CASCADE
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE operation_vars_map;
DROP TABLE operation_vars;
