// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"strings"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type EditEvidenceMetadataInput struct {
	OperationSlug string
	EvidenceUUID  string
	Source        string
	Body          string
}

type UpsertEvidenceMetadataInput struct {
	EditEvidenceMetadataInput
	Status     string
	Message    *string
	CanProcess *bool
}

func CreateEvidenceMetadata(ctx context.Context, db *database.Connection, i EditEvidenceMetadataInput) error {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.WrapError("Unable to create evidence metadata", backend.UnauthorizedWriteErr(err))
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to create evidence metadata", backend.UnauthorizedWriteErr(err))
	}

	_, err = db.Insert("evidence_metadata", map[string]interface{}{
		"evidence_id": evidence.ID,
		"source":      i.Source,
		"body":        i.Body,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return backend.WrapError("Couold not edit evidence metadata",
				backend.SuggestiveDatabaseErr("This metadata source already exists", err),
			)
		}

		return backend.WrapError("Could not create evidence metadata", backend.DatabaseErr(err))
	}

	return nil
}

func UpdateEvidenceMetadata(ctx context.Context, db *database.Connection, i EditEvidenceMetadataInput) error {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.WrapError("Unable to edit evidence metadata", backend.UnauthorizedWriteErr(err))
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to edit evidence metadata", backend.UnauthorizedWriteErr(err))
	}

	err = db.Update(sq.
		Update("evidence_metadata").
		Set("body", i.Body).
		Where(sq.Eq{
			"evidence_id": evidence.ID,
			"source":      i.Source,
		}))

	if err != nil {
		return backend.WrapError("Could not edit evidence metadata", backend.DatabaseErr(err))
	}

	return nil
}

func UpsertEvidenceMetadata(ctx context.Context, db *database.Connection, i UpsertEvidenceMetadataInput) error {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.WrapError("Unable to edit evidence metadata", backend.UnauthorizedWriteErr(err))
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to edit evidence metadata", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		var metadata []models.EvidenceMetadata
		tx.Select(&metadata, sq.Select("*").From("evidence_metadata").Where(sq.Eq{
			"evidence_id": evidence.ID,
			"source":      i.Source,
		}))
		// these should call out to helper functions to do the work (shared with the true methods),
		// but we need some db work to be integrated before we can do this properly.
		if len(metadata) == 0 {
			tx.Insert("evidence_metadata", map[string]interface{}{
				"evidence_id":      evidence.ID,
				"source":           i.Source,
				"body":             i.Body,
				"last_run_message": i.Message,
				"can_process":      i.CanProcess,
			})
		} else {
			tx.Update(sq.
				Update("evidence_metadata").
				SetMap(map[string]interface{}{
					"body":             i.Body,
					"last_run_message": i.Message,
					"can_process":      i.CanProcess,
					"status":           i.Status,
				}).
				Where(sq.Eq{
					"evidence_id": evidence.ID,
					"source":      i.Source,
				}))
		}
	})
	if err != nil {
		return backend.WrapError("Could not edit evidence metadata", backend.DatabaseErr(err))
	}

	return nil
}
