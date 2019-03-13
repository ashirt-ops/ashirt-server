// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

// listAllOperations is a helper function for both ListOperations and ListOpperationsForAdmin.
// This retrieves all operations, then relies on the caller to sort which operations are visible
// to the enduser
func listAllOperations(db *database.Connection) ([]*dtos.Operation, error) {
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
		return nil, backend.DatabaseErr(err)
	}

	operationsDTO := []*dtos.Operation{}
	for _, operation := range operations {
		operationsDTO = append(operationsDTO, &dtos.Operation{
			Slug:     operation.Slug,
			Name:     operation.Name,
			Status:   operation.Status,
			NumUsers: operation.NumUsers,

			// Temporary for screenshot client:
			ID: operation.ID,
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
			operationsDTO = append(operationsDTO, operation)
		}
	}
	return operationsDTO, nil
}
