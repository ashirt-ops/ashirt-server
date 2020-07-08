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

// DeleteUser needs some godocs
func DeleteUser(ctx context.Context, db *database.Connection, slug string) error {
	if !middleware.IsAdmin(ctx) {
		return backend.UnauthorizedWriteErr(fmt.Errorf("Requesting user is not an admin"))
	}

	userID, err := userSlugToUserID(db, slug)
	if err != nil {
		return backend.DatabaseErr(err)
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
		return disableErr
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("api_keys").Where(sq.Eq{"user_id": userID}))
		tx.Delete(sq.Delete("auth_scheme_data").Where(sq.Eq{"user_id": userID}))
		tx.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"user_id": userID}))
		tx.Update(sq.Update("users").Set("deleted_at", time.Now()).Where(sq.Eq{"slug": slug}))
	})
	if err != nil {
		return backend.DatabaseErr(err)
	}

	return nil
}
