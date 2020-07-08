// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ReadFindingInput struct {
	OperationSlug string
	FindingUUID   string
}

func ReadFinding(ctx context.Context, db *database.Connection, i ReadFindingInput) (*dtos.Finding, error) {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var evidenceIDs []int64

	err = db.Select(&evidenceIDs, sq.Select("evidence_id").
		From("evidence_finding_map").
		Where(sq.Eq{"finding_id": finding.ID}))
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	_, allTags, err := tagsForEvidenceByID(db, evidenceIDs)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	return &dtos.Finding{
		UUID:          i.FindingUUID,
		Title:         finding.Title,
		Category:      finding.Category,
		Description:   finding.Description,
		NumEvidence:   len(evidenceIDs),
		Tags:          allTags,
		ReadyToReport: finding.ReadyToReport,
		TicketLink:    finding.TicketLink,
	}, nil
}
