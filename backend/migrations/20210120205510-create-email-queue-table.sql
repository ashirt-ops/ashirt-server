-- +migrate Up
CREATE TABLE email_queue (
  id INT AUTO_INCREMENT,
  to_email VARCHAR(255) NOT NULL,
  user_id INT NOT NULL DEFAULT 0,
  template VARCHAR(255) NOT NULL,
  email_sent BOOLEAN NOT NULL DEFAULT false,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
CREATE INDEX email_queue__email_sent ON email_queue(email_sent);
CREATE INDEX email_queue__email_to ON email_queue(to_email);

ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE(email);

-- +migrate Down
ALTER TABLE users DROP CONSTRAINT unique_email;
DROP INDEX email_queue__email_to on email_queue;
DROP INDEX email_queue__email_sent on email_queue;
DROP TABLE email_queue;
