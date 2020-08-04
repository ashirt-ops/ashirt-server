// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

// DeleteUser provides the ability for a super admin to remove a user from the system. Doing so
// removes access only. Evidence and other contributions remain. Note that users are not able to
// delete their own accounts to prevent accidents. Also note that once a user has been deleted,
// they cannot be restored.
func DeleteUser(ctx context.Context, db *database.Connection, slug string) error {
	if !middleware.IsAdmin(ctx) {
		return backend.WrapError("Unwilling to delete user", backend.UnauthorizedWriteErr(fmt.Errorf("Requesting user is not an admin")))
	}

	userID, err := userSlugToUserID(db, slug)
	if err != nil {
		return backend.WrapError("Unable to delete user", backend.DatabaseErr(err))
	}

	if userID == middleware.UserID(ctx) {
		return backend.BadInputErr(fmt.Errorf("User is trying to delete themself"), "Users cannot delete themselves")
	}

	disabled := true
	// session data is deleted when disabling the user
	disableErr := SetUserFlags(ctx, db, SetUserFlagsInput{
		Slug:     slug,
		Disabled: &disabled,
	})
	if disableErr != nil {
		return backend.WrapError("Could not set user to disabled prior to deletion", disableErr)
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("api_keys").Where(sq.Eq{"user_id": userID}))
		tx.Delete(sq.Delete("auth_scheme_data").Where(sq.Eq{"user_id": userID}))
		tx.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"user_id": userID}))
		tx.Update(sq.Update("users").Set("deleted_at", time.Now()).Where(sq.Eq{"slug": slug}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete user", backend.DatabaseErr(err))
	}

	return nil
}
