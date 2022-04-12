// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	// sq "github.com/Masterminds/squirrel"
)

type CreateEvidenceMetadataInput struct {
	OperationSlug string
	EvidenceUUID  string
	Source        string
	Body          string
}

func CreateEvidenceMetadata(ctx context.Context, db *database.Connection, i CreateEvidenceMetadataInput) error {
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
		return backend.WrapError("Could not create evidence metadata", backend.DatabaseErr(err))
	}

	return nil
}
