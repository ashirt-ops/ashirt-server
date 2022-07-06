// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices

import (
	"context"
	"fmt"

	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"

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

func getExpandedPayloadID(e ExpandedNewEvidencePayload) int64 {
	return e.EvidenceID
}

// BatchBuildNewEvidencePayload creates a set of payloads for the given operation and evidence uuids.
// This function provides convenience over the alternatives: BatchBuildNewEvidencePayloadFromUUIDs and
// BatchBuildNewEvidencePayloadForAllEvidence. Note that if no evidenceUUIDs are provided,
// then all evidence is chosen for the indicated operation.
func BatchBuildNewEvidencePayload(ctx context.Context, db database.ConnectionProxy, operationID int64, evidenceUUIDs []string) ([]ExpandedNewEvidencePayload, error) {
	if len(evidenceUUIDs) == 0 {
		return BatchBuildNewEvidencePayloadForAllEvidence(ctx, db, operationID)
	} else {
		return BatchBuildNewEvidencePayloadFromUUIDs(ctx, db, operationID, evidenceUUIDs)
	}
}

// BatchBuildNewEvidencePayloadFromUUIDs creates a set of payloads, ordered by evidence ID, for the given operationID and evidenceUUIDs.
// If the list of uuids is empty, then _no payloads will be returned_.
// Also see BatchBuildNewEvidencePayload, which allows for getting all evidence if no uuids are specified
func BatchBuildNewEvidencePayloadFromUUIDs(ctx context.Context, db database.ConnectionProxy, operationID int64, evidenceUUIDs []string) ([]ExpandedNewEvidencePayload, error) {
	return batchBuildNewEvidencePayloadSpecial(ctx, db, func(tx database.ConnectionProxy) ([]models.Evidence, error) {
		return database.GetEvidenceFromUUIDs(tx, operationID, evidenceUUIDs)
	})
}

// BatchBuildNewEvidencePayloadForAllEvidence creates a set of payloads, ordered by evidence ID, for all evidence in an operation.
// Also see BatchBuildNewEvidencePayload, which allows for specifying a subset of evidence uuids
func BatchBuildNewEvidencePayloadForAllEvidence(ctx context.Context, db database.ConnectionProxy, operationID int64) ([]ExpandedNewEvidencePayload, error) {
	return batchBuildNewEvidencePayloadSpecial(ctx, db, func(tx database.ConnectionProxy) ([]models.Evidence, error) {
		return database.GetAllEvidenceForOperation(tx, operationID)
	})
}

func batchBuildNewEvidencePayloadSpecial(ctx context.Context, db database.ConnectionProxy,
	fetch func(tx database.ConnectionProxy) ([]models.Evidence, error),
) ([]ExpandedNewEvidencePayload, error) {
	var payloads []ExpandedNewEvidencePayload
	err := db.WithTx(ctx, func(tx *database.Transactable) {
		evidence, _ := fetch(tx)
		ids := helpers.Map(evidence, database.EvidenceToID)
		payloads, _ = batchBuildNewEvidencePayloadFromIDs(tx, ids)
	})

	return payloads, err
}

// batchBuildNewEvidencePayloadFromIDs builds a payload by getting all of the necessary details in bulk.
// Note: this relies on the ordering of evidenceIDs. No particular order is required as input,
// but the result is ordered by evidenceID, in ASC order.
func batchBuildNewEvidencePayloadFromIDs(db database.ConnectionProxy, evidenceIDs []int64) ([]ExpandedNewEvidencePayload, error) {
	var payloads []ExpandedNewEvidencePayload

	err := db.Select(&payloads, sq.Select(
		"e.id AS id",
		"e.uuid AS uuid",
		"e.content_type",
		"slug AS operation_slug",
		"'evidence_created' AS type", // hardcode in the type so we don't have to edit each entry manually
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
