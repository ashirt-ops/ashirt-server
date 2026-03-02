-- +migrate Up
CREATE TABLE evidence_event_map (
  evidence_id INT NOT NULL,
  event_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (evidence_id, event_id),
  FOREIGN KEY (evidence_id) REFERENCES evidence(id) ON DELETE CASCADE,
  FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE evidence_event_map;
