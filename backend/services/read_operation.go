// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"

	sq "github.com/Masterminds/squirrel"
)

func ReadOperation(ctx context.Context, db *database.Connection, operationSlug string) (*dtos.Operation, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var numUsers int
	err = db.Get(&numUsers, sq.Select("count(*)").From("user_operation_permissions").
		Where(sq.Eq{"operation_id": operation.ID}))
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	return &dtos.Operation{
		Slug:     operationSlug,
		Name:     operation.Name,
		Status:   operation.Status,
		NumUsers: numUsers,
	}, nil
}
