package services

import (
	"context"
	"fmt"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"

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
		return errors.WrapError("Unable to set user role", errors.UnauthorizedWriteErr(err))
	}

	if i.UserSlug == "" {
		return errors.MissingValueErr("User Slug")
	}

	userID, err := userSlugToUserID(db, i.UserSlug)
	if err != nil {
		return errors.WrapError("Unable to get user id from slug", errors.BadInputErr(err, fmt.Sprintf(`No user with slug "%s" was found`, i.UserSlug)))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyUserOfOperation{UserID: userID, OperationID: operation.ID}); err != nil {
		return errors.WrapError("Unwilling to set user role", errors.UnauthorizedWriteErr(err))
	}

	if i.Role == "" {
		err := db.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"user_id": userID, "operation_id": operation.ID}))

		if err != nil {
			return errors.WrapError("Cannot delete user role", errors.DatabaseErr(err))
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
			return errors.WrapError("Unable to add user role", errors.DatabaseErr(err))
		}
		return nil
	}

	if permission.Role != i.Role {
		err = db.Update(sq.Update("user_operation_permissions").
			Set("role", i.Role).
			Where(sq.Eq{"user_id": userID, "operation_id": operation.ID}))

		if err != nil {
			return errors.WrapError("Unable to alter user role", errors.DatabaseErr(err))
		}
	}
	return nil
}

func SetUserGroupOperationRole(ctx context.Context, db *database.Connection, i SetUserGroupOperationRoleInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return errors.WrapError("Unable to set user group role", errors.UnauthorizedWriteErr(err))
	}

	if i.UserGroupSlug == "" {
		return errors.MissingValueErr("User Group Slug")
	}

	userGroupID, err := userGroupSlugToUserGroupID(db, i.UserGroupSlug)
	if err != nil {
		return errors.WrapError("Unable to get user group id from slug", errors.BadInputErr(err, fmt.Sprintf(`No user with slug "%s" was found`, i.UserGroupSlug)))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyUserGroupOfOperation{UserGroupID: userGroupID, OperationID: operation.ID}); err != nil {
		return errors.WrapError("Unwilling to set user group role", errors.UnauthorizedWriteErr(err))
	}

	if i.Role == "" {
		err := db.Delete(sq.Delete("user_group_operation_permissions").Where(sq.Eq{"group_id": userGroupID, "operation_id": operation.ID}))

		if err != nil {
			return errors.WrapError("Cannot delete user group role", errors.DatabaseErr(err))
		}
		return nil
	}

	var permissions []models.UserGroupOperationPermission
	err = db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.Select(&permissions, sq.Select("*").
			From("user_group_operation_permissions").
			Where(sq.Eq{
				"group_id":     userGroupID,
				"operation_id": operation.ID,
			}))
		if len(permissions) == 0 {
			tx.Insert("user_group_operation_permissions", map[string]interface{}{
				"group_id":     userGroupID,
				"operation_id": operation.ID,
				"role":         i.Role,
			})
		} else if permissions[0].Role != i.Role {
			tx.Update(sq.Update("user_group_operation_permissions").
				Set("role", i.Role).
				Where(sq.Eq{"group_id": userGroupID, "operation_id": operation.ID}))
		}
	})
	if err != nil {
		return errors.WrapError("Unable to add user role", errors.DatabaseErr(err))
	}

	return nil
}
