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

type UpdateFindingInput struct {
	OperationSlug string
	FindingUUID   string
	Category      string
	Title         string
	Description   string
	TicketLink    *string
	ReadyToReport bool
}

func UpdateFinding(ctx context.Context, db *database.Connection, i UpdateFindingInput) error {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return backend.WrapError("Unable to lookup operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Failed permission check", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		useCategoryID, _ := getFindingCategoryID(i.Category, tx.Select)

		tx.Update(sq.Update("findings").
			SetMap(map[string]interface{}{
				"category_id":     useCategoryID,
				"title":           i.Title,
				"description":     i.Description,
				"ticket_link":     i.TicketLink,
				"ready_to_report": i.ReadyToReport,
			}).
			Where(sq.Eq{"id": finding.ID}))
	})

	if err != nil {
		return backend.WrapError("Unable to update database", backend.UnauthorizedWriteErr(err))
	}
	return nil
}
