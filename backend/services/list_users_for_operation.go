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

	sq "github.com/Masterminds/squirrel"
)

type ListUsersForOperationInput struct {
	Pagination
	UserFilter
	OperationSlug string
}

type userAndRole struct {
	models.User
	Role policy.OperationRole `db:"role"`
}

func ListUsersForOperation(ctx context.Context, db *database.Connection, i ListUsersForOperationInput) ([]*dtos.UserOperationRole, error) {
	query, err := prepListUsersForOperation(ctx, db, i)
	if err != nil {
		return nil, err
	}

	var users []userAndRole
	err = db.Select(&users, *query)
	if err != nil {
		return nil, backend.WrapError("Cannot list users for operation", backend.DatabaseErr(err))
	}
	usersDTO := wrapListUsersForOperationResponse(users)
	return usersDTO, nil
}

func prepListUsersForOperation(ctx context.Context, db *database.Connection, i ListUsersForOperationInput) (*sq.SelectBuilder, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list users for operation", backend.UnauthorizedReadErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanListUsersOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list users for operation", backend.UnauthorizedReadErr(err))
	}

	query := sq.Select("slug", "first_name", "last_name", "role").
		From("user_operation_permissions").
		LeftJoin("users ON user_operation_permissions.user_id = users.id").
		Where(sq.Eq{"operation_id": operation.ID, "users.deleted_at": nil}).
		OrderBy("user_operation_permissions.created_at ASC")
		// OrderBy("first_name ASC", "last_name ASC", "user_operation_permissions.created_at ASC")

	i.UserFilter.AddWhere(&query)

	return &query, nil
}

func wrapListUsersForOperationResponse(users []userAndRole) []*dtos.UserOperationRole {
	usersDTO := make([]*dtos.UserOperationRole, len(users))
	for idx, user := range users {
		usersDTO[idx] = &dtos.UserOperationRole{
			User: dtos.User{
				Slug:      user.Slug,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
			Role: user.Role,
		}
	}
	return usersDTO
}
