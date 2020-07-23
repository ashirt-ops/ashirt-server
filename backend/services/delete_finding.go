// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type DeleteFindingInput struct {
	OperationSlug string
	FindingUUID   string
}

func DeleteFinding(ctx context.Context, db *database.Connection, i DeleteFindingInput) error {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return backend.WrapError("Unable to delete finding", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete finding", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"finding_id": finding.ID}))
		tx.Delete(sq.Delete("findings").Where(sq.Eq{"id": finding.ID}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete finding", backend.DatabaseErr(err))
	}

	return nil
}
