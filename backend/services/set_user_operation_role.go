// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"

	sq "github.com/Masterminds/squirrel"
)

type SetUserOperationRoleInput struct {
	OperationSlug string
	UserSlug      string
	Role          policy.OperationRole
}

func SetUserOperationRole(ctx context.Context, db *database.Connection, i SetUserOperationRoleInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	if i.UserSlug == "" {
		return backend.MissingValueErr("User Slug")
	}

	userID, err := userSlugToUserID(db, i.UserSlug)
	if err != nil {
		return backend.BadInputErr(err, fmt.Sprintf(`No user with slug "%s" was found`, i.UserSlug))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyUserOfOperation{UserID: userID, OperationID: operation.ID}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	if i.Role == "" {
		err := db.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"user_id": userID, "operation_id": operation.ID}))

		if err != nil {
			return backend.DatabaseErr(err)
		}
		return nil
	}

	var permission models.UserOperationPermission
	err = db.Get(&permission, sq.Select("*").
		From("user_operation_permissions").
		Where(sq.Eq{
			"user_id":      userID,
			"operation_id": operation.ID,
		}))
	if err != nil {
		_, err = db.Insert("user_operation_permissions", map[string]interface{}{
			"user_id":      userID,
			"operation_id": operation.ID,
			"role":         i.Role,
		})
		if err != nil {
			return backend.DatabaseErr(err)
		}
		return nil
	}

	if permission.Role != i.Role {
		err = db.Update(sq.Update("user_operation_permissions").
			Set("role", i.Role).
			Where(sq.Eq{"user_id": userID, "operation_id": operation.ID}))

		if err != nil {
			return backend.DatabaseErr(err)
		}
	}
	return nil
}
