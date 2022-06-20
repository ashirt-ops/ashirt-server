package enhancementservices

import (
	"fmt"

	"github.com/theparanoids/ashirt-server/backend/database"

	sq "github.com/Masterminds/squirrel"
)

type NewEvidencePayload struct {
	Type          string `json:"type" db:"type"`
	EvidenceUUID  string `json:"evidenceUuid"  db:"uuid"`
	OperationSlug string `json:"operationSlug" db:"operation_slug"`
	ContentType   string `json:"contentType"   db:"content_type"`
}

type ExpandedNewEvidencePayload struct {
	NewEvidencePayload
	EvidenceID int64 `db:"id"`
}

// BatchBuildNewEvidencePayload builds a payload by getting all of the necessary details in bulk.
// Note: this relies on the ordering of evidenceIDs. No particular order is required as input,
// but the result is ordered by evidenceID, in ASC order.
func BatchBuildNewEvidencePayload(db database.ConnectionProxy, evidenceIDs []int64) ([]ExpandedNewEvidencePayload, error) {
	var payloads []ExpandedNewEvidencePayload

	err := db.Select(&payloads, sq.Select(
		"e.id AS id",
		"e.uuid AS uuid",
		"e.content_type",
		"slug AS operation_slug",
		"'process' AS type", // hardcode in the type so we don't have to edit each entry manually
	).
		From("evidence e").
		LeftJoin("operations o ON e.operation_id = o.id").
		Where(sq.Eq{"e.id": evidenceIDs}).
		OrderBy(`e.id`),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to gather evidence data for worker")
	}

	return payloads, nil
}
