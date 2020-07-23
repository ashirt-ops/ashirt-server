// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"

	sq "github.com/Masterminds/squirrel"
)

type UpdateOperationInput struct {
	OperationSlug string
	Name          string
	Status        models.OperationStatus
}

func UpdateOperation(ctx context.Context, db *database.Connection, i UpdateOperationInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to update operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to update operation", backend.UnauthorizedWriteErr(err))
	}

	err = db.Update(sq.Update("operations").
		SetMap(map[string]interface{}{
			"name":   i.Name,
			"status": i.Status,
		}).
		Where(sq.Eq{"id": operation.ID}))
	if err != nil {
		return backend.WrapError("Cannot update operation", backend.DatabaseErr(err))
	}
	return nil
}
