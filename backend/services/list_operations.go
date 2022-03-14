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

type operationListItem struct {
	Op *dtos.Operation
	ID int64
}

// listAllOperations is a helper function for both ListOperations and ListOpperationsForAdmin.
// This retrieves all operations, then relies on the caller to sort which operations are visible
// to the enduser
func listAllOperations(db *database.Connection) ([]operationListItem, error) {
	var operations []struct {
		models.Operation
		NumUsers int `db:"num_users"`
	}

	err := db.Select(&operations, sq.Select("id", "slug", "name", "status", "count(user_id) AS num_users").
		From("operations").
		LeftJoin("user_operation_permissions ON user_operation_permissions.operation_id = operations.id").
		GroupBy("operations.id").
		OrderBy("operations.created_at DESC"))
	if err != nil {
		return nil, backend.WrapError("Cannot list all operations", backend.DatabaseErr(err))
	}

	operationsDTO := []operationListItem{}
	for _, operation := range operations {
		operationsDTO = append(operationsDTO, operationListItem{
			ID: operation.ID,
			Op: &dtos.Operation{
				Slug:     operation.Slug,
				Name:     operation.Name,
				Status:   operation.Status,
				NumUsers: operation.NumUsers,
			},
		})
	}
	return operationsDTO, nil
}

// ListOperations retrieves a list of all operations that the contextual user can see
func ListOperations(ctx context.Context, db *database.Connection) ([]*dtos.Operation, error) {
	operations, err := listAllOperations(db)

	if err != nil {
		return nil, err
	}

	operationsDTO := make([]*dtos.Operation, 0, len(operations))
	for _, operation := range operations {
		if middleware.Policy(ctx).Check(policy.CanReadOperation{OperationID: operation.ID}) {
			operationsDTO = append(operationsDTO, operation.Op)
		}
	}
	return operationsDTO, nil
}
