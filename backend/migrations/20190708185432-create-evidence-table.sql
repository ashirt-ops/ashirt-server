-- +migrate Up
CREATE TABLE evidence (
  id INT AUTO_INCREMENT,
  operation_id INT NOT NULL,
  operator_id INT NOT NULL,
  description TEXT,
  full_image_key VARCHAR(255),
  thumb_image_key VARCHAR(255),
  occurred_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (operation_id) REFERENCES operations(id),
  FOREIGN KEY (operator_id) REFERENCES users(id)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- +migrate Down
DROP TABLE evidence;
