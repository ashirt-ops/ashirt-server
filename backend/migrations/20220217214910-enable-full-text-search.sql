-- +migrate Up
CREATE TABLE `evidence_metadata` (
    id INT AUTO_INCREMENT,
    evidence_id INT NOT NULL,
    source varchar(255) NOT NULL,
    metadata TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY `fk_evidence_id` (evidence_id) REFERENCES evidence(id),
    FULLTEXT `i__metadata` (metadata)
) ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;

-- +migrate Down
DROP TABLE `evidence_metadata`;
