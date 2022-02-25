// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type UpdateEvidenceMetadataInput struct {
	MetadataId    int64
	OperationSlug string
	EvidenceUUID  string
	Source        string
	Metadata      string
}

func UpdateEvidenceMetadata(
	ctx context.Context,
	db *database.Connection,
	contentStore contentstore.Store,
	i UpdateEvidenceMetadataInput,
) error {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.WrapError("Unable to update evidence", backend.UnauthorizedWriteErr(err))
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to update evidence", backend.UnauthorizedWriteErr(err))
	}

	err = db.Update(sq.Update("evidence_metadata").
		SetMap(map[string]interface{}{
			"metadta": i.Metadata,
		}).
		Where(sq.Eq{"evidence_id": evidence.ID}))

	if err != nil {
		return backend.WrapError("Could not update evidence metadata", backend.DatabaseErr(err))
	}

	return nil
}
