-- +migrate Up
CREATE TABLE default_tags (
    id INT AUTO_INCREMENT,
    name VARCHAR(63) NOT NULL UNIQUE,
    color_name VARCHAR(63) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;

-- +migrate Down
DROP TABLE default_tags;
