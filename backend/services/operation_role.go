// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"

	sq "github.com/Masterminds/squirrel"
)

type SetUserOperationRoleInput struct {
	OperationSlug string
	UserSlug      string
	Role          policy.OperationRole
}

type SetUserGroupOperationRoleInput struct {
	OperationSlug string
	UserGroupSlug string
	Role          policy.OperationRole
}

func SetUserOperationRole(ctx context.Context, db *database.Connection, i SetUserOperationRoleInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to set user role", backend.UnauthorizedWriteErr(err))
	}

	if i.UserSlug == "" {
		return backend.MissingValueErr("User Slug")
	}

	userGroupID, err := userSlugToUserID(db, i.UserSlug)
	if err != nil {
		return backend.WrapError("Unable to get user id from slug", backend.BadInputErr(err, fmt.Sprintf(`No user with slug "%s" was found`, i.UserSlug)))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyUserOfOperation{UserID: userGroupID, OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to set user role", backend.UnauthorizedWriteErr(err))
	}

	if i.Role == "" {
		err := db.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"user_id": userGroupID, "operation_id": operation.ID}))

		if err != nil {
			return backend.WrapError("Cannot delete user role", backend.DatabaseErr(err))
		}
		return nil
	}

	var permission models.UserOperationPermission
	err = db.Get(&permission, sq.Select("*").
		From("user_operation_permissions").
		Where(sq.Eq{
			"user_id":      userGroupID,
			"operation_id": operation.ID,
		}))
	if err != nil {
		_, err = db.Insert("user_operation_permissions", map[string]interface{}{
			"user_id":      userGroupID,
			"operation_id": operation.ID,
			"role":         i.Role,
		})
		if err != nil {
			return backend.WrapError("Unable to add user role", backend.DatabaseErr(err))
		}
		return nil
	}

	if permission.Role != i.Role {
		err = db.Update(sq.Update("user_operation_permissions").
			Set("role", i.Role).
			Where(sq.Eq{"user_id": userGroupID, "operation_id": operation.ID}))

		if err != nil {
			return backend.WrapError("Unable to alter user role", backend.DatabaseErr(err))
		}
	}
	return nil
}

func SetUserGroupOperationRole(ctx context.Context, db *database.Connection, i SetUserGroupOperationRoleInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to set user group role", backend.UnauthorizedWriteErr(err))
	}

	if i.UserGroupSlug == "" {
		return backend.MissingValueErr("User Group Slug")
	}

	userGroupID, err := userGroupSlugToUserGroupID(db, i.UserGroupSlug)
	if err != nil {
		return backend.WrapError("Unable to get user group id from slug", backend.BadInputErr(err, fmt.Sprintf(`No user with slug "%s" was found`, i.UserGroupSlug)))
	}

	// TODO TN create policy
	// if err := policyRequireWithAdminBypass(ctx, policy.CanModifyUserOfOperation{UserID: userGroupID, OperationID: operation.ID}); err != nil {
	// 	return backend.WrapError("Unwilling to set user group role", backend.UnauthorizedWriteErr(err))
	// }

	if i.Role == "" {
		err := db.Delete(sq.Delete("user_group_operation_permissions").Where(sq.Eq{"group_id": userGroupID, "operation_id": operation.ID}))

		if err != nil {
			return backend.WrapError("Cannot delete user group role", backend.DatabaseErr(err))
		}
		return nil
	}

	var permission models.UserGroupOperationPermission
	err = db.Get(&permission, sq.Select("*").
		From("user_group_operation_permissions").
		Where(sq.Eq{
			"group_id":     userGroupID,
			"operation_id": operation.ID,
		}))
	if err != nil {
		_, err = db.Insert("user_group_operation_permissions", map[string]interface{}{
			"group_id":     userGroupID,
			"operation_id": operation.ID,
			"role":         i.Role,
		})
		if err != nil {
			return backend.WrapError("Unable to add user role", backend.DatabaseErr(err))
		}
		return nil
	}

	if permission.Role != i.Role {
		err = db.Update(sq.Update("user_group_operation_permissions").
			Set("role", i.Role).
			Where(sq.Eq{"group_id": userGroupID, "operation_id": operation.ID}))

		if err != nil {
			return backend.WrapError("Unable to alter user role", backend.DatabaseErr(err))
		}
	}
	return nil
}
