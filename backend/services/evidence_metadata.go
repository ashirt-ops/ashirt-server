// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateEvidenceMetadataInput struct {
	OperationSlug string
	EvidenceUUID  string
	Source        string
	Metadata      string
}

type UpdateEvidenceMetadataInput struct {
	MetadataId    int64
	OperationSlug string
	EvidenceUUID  string
	Source        string
	Metadata      string
}

func CreateEvidenceMetadata(
	ctx context.Context,
	db *database.Connection,
	contentStore contentstore.Store,
	i CreateEvidenceMetadataInput,
) (*dtos.EvidenceMetadata, error) {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return nil, backend.WrapError("Unable to update evidence", backend.UnauthorizedWriteErr(err))
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to update evidence", backend.UnauthorizedWriteErr(err))
	}

	insertID, err := db.Insert("evidence_metadata", map[string]interface{}{
		"evidenceId": evidence.ID,
		"source":     i.Source,
		"metadta":    i.Metadata,
	})

	if err != nil {
		return nil, backend.WrapError("Could not create evidence metadata", backend.DatabaseErr(err))
	}

	return &dtos.EvidenceMetadata{
		ID:       insertID,
		Metadata: i.Metadata,
		Source:   i.Source,
	}, nil
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
