
-- +migrate Up
ALTER TABLE evidence
	ADD COLUMN adjusted_at TIMESTAMP
;

-- +migrate Down
ALTER TABLE evidence 
	DROP COLUMN adjusted_at
;