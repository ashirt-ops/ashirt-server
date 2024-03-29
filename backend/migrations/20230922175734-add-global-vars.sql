-- +migrate Up
CREATE TABLE global_vars (
  id INT AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,
  value VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE global_vars;
