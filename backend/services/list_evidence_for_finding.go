// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ListEvidenceForFindingInput struct {
	OperationSlug string
	FindingUUID   string
}

func ListEvidenceForFinding(ctx context.Context, db *database.Connection, i ListEvidenceForFindingInput) ([]dtos.Evidence, error) {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return nil, backend.WrapError("Unable to list evidence for finding", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list evidence for finding", backend.UnauthorizedReadErr(err))
	}

	var evidenceForFinding []struct {
		models.Evidence
		Slug      string `db:"slug"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	err = db.Select(&evidenceForFinding, sq.Select("evidence.*", "slug", "first_name", "last_name").
		From("evidence").
		LeftJoin("evidence_finding_map ON evidence.id = evidence_id").
		LeftJoin("users ON users.id = evidence.operator_id").
		Where(sq.Eq{"finding_id": finding.ID}))

	if err != nil {
		return nil, backend.WrapError("Cannot list evidence for finding", backend.UnauthorizedReadErr(err))
	}

	evidenceIDs := make([]int64, len(evidenceForFinding))
	for idx, evi := range evidenceForFinding {
		evidenceIDs[idx] = evi.Evidence.ID
	}

	tagsByEvidenceID, _, err := tagsForEvidenceByID(db, evidenceIDs)
	if err != nil {
		return nil, backend.WrapError("Cannot get tags for evidnece", backend.UnauthorizedReadErr(err))
	}

	var evidenceDTOs = make([]dtos.Evidence, len(evidenceForFinding))
	for i, evi := range evidenceForFinding {
		tags := tagsByEvidenceID[evi.Evidence.ID]
		if tags == nil {
			tags = []dtos.Tag{}
		}
		evidenceDTOs[i] = dtos.Evidence{
			UUID:        evi.UUID,
			ContentType: evi.ContentType,
			Description: evi.Description,
			OccurredAt:  evi.OccurredAt,
			Tags:        tags,
			Operator: dtos.User{
				Slug:      evi.Slug,
				FirstName: evi.FirstName,
				LastName:  evi.LastName,
			},
		}
	}

	return evidenceDTOs, nil
}
