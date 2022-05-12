// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/sync/errgroup"
)

func DeleteOperation(ctx context.Context, db *database.Connection, contentStore contentstore.Store, slug string) error {
	operation, err := lookupOperation(db, slug)
	if err != nil {
		return backend.WrapError("Unable to delete operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanDeleteOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete operation", backend.UnauthorizedWriteErr(err))
	}
	log := logging.ReqLogger(ctx)

	var g errgroup.Group
	g.Go(func() error {
		err := db.WithTx(ctx, func(tx *database.Transactable) {
			var evidence []models.Evidence
			err = tx.Select(&evidence, sq.Select("*").From("evidence").Where(sq.Eq{"operation_id": operation.ID}))

			// remove evidence content
			if err == nil {
				for _, evi := range evidence {
					copy := evi
					g.Go(func() error {
						err := deleteEvidenceContent(contentStore, copy)
						if err != nil {
							log.Log("task", "delete operation", "msg", "error deleting evidence content", "uniqueKey", "orphanedDelete",
								"keys", fmt.Sprintf(`["%v", "%v"]`, copy.FullImageKey, copy.ThumbImageKey), "error", err.Error())
							return backend.DeleteErr(err)
						}
						return nil
					})
				}
			}

			// remove all tags for an operation
			var tagIDs []int64
			tx.Select(&tagIDs, sq.Select("id").From("tags").Where(sq.Eq{"operation_id": operation.ID}))
			tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"tag_id": tagIDs}))
			tx.Delete(sq.Delete("tags").Where(sq.Eq{"id": tagIDs}))

			// remove all findings for an operation
			var findingIDs []int64
			tx.Select(&findingIDs, sq.Select("id").From("findings").Where(sq.Eq{"operation_id": operation.ID}))
			tx.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"finding_id": findingIDs}))
			tx.Delete(sq.Delete("findings").Where(sq.Eq{"id": findingIDs}))

			var evidenceIDs = make([]int64, len(evidence))
			for i, evi := range evidence {
				evidenceIDs[i] = evi.ID
			}

			// remove evidence metadata
			tx.Delete(sq.Delete("evidence_metadata").Where(sq.Eq{"evidence_id": evidenceIDs}))

			// remove all evidence
			tx.Delete(sq.Delete("evidence").Where(sq.Eq{"id": evidenceIDs}))

			// remove user/operations map
			tx.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"operation_id": operation.ID}))

			tx.Delete(sq.Delete("operations").Where(sq.Eq{"id": operation.ID}))
		})
		if err != nil {
			log.Log("task", "delete operation", "msg", "Failed to fully delete operation data",
				"error", err.Error())
			return backend.WrapError("Cannot delete operation", backend.DatabaseErr(err))
		}
		return nil
	})

	return g.Wait()
}
