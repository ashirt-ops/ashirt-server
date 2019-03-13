// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
	"golang.org/x/sync/errgroup"

	sq "github.com/Masterminds/squirrel"
)

type AddEvidenceToFindingInput struct {
	OperationSlug    string
	FindingUUID      string
	EvidenceToAdd    []string
	EvidenceToRemove []string
}

func AddEvidenceToFinding(ctx context.Context, db *database.Connection, i AddEvidenceToFindingInput) error {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	var g errgroup.Group
	g.Go(func() (err error) { return batchAddEvidenceToFinding(db, i.EvidenceToAdd, operation.ID, finding.ID) })
	g.Go(func() (err error) { return batchRemoveEvidenceFromFinding(db, i.EvidenceToRemove, finding.ID) })
	if err = g.Wait(); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	return nil
}

func buildQueryForEvidenceFromUUIDs(evidenceUUIDs []string) sq.SelectBuilder {
	return sq.Select("*").
		From("evidence").
		Where(sq.Eq{"uuid": evidenceUUIDs})
}

func batchAddEvidenceToFinding(db *database.Connection, evidenceUUIDs []string, operationID int64, findingID int64) error {
	if len(evidenceUUIDs) == 0 {
		return nil
	}
	var evidence []models.Evidence
	if err := db.Select(&evidence, buildQueryForEvidenceFromUUIDs(evidenceUUIDs)); err != nil {
		return err
	}
	evidenceIDs := []int64{}
	for _, evi := range evidence {
		if evi.OperationID != operationID {
			return fmt.Errorf(
				"Cannot add evidence %d to operation %d. Evidence belongs to operation %d",
				evi.ID, operationID, evi.OperationID,
			)
		}
		evidenceIDs = append(evidenceIDs, evi.ID)
	}
	return db.BatchInsert("evidence_finding_map", len(evidenceIDs), func(idx int) map[string]interface{} {
		return map[string]interface{}{
			"finding_id":  findingID,
			"evidence_id": evidenceIDs[idx],
		}
	})
}

func batchRemoveEvidenceFromFinding(db *database.Connection, evidenceUUIDs []string, findingID int64) error {
	if len(evidenceUUIDs) == 0 {
		return nil
	}
	var evidence []models.Evidence
	if err := db.Select(&evidence, buildQueryForEvidenceFromUUIDs(evidenceUUIDs)); err != nil {
		return err
	}
	evidenceIDs := []int64{}
	for _, evi := range evidence {
		evidenceIDs = append(evidenceIDs, evi.ID)
	}

	return db.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"finding_id": findingID, "evidence_id": evidenceIDs}))
}
