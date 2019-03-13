// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/server/middleware"
	sq "github.com/Masterminds/squirrel"
)

type SetUserFlagsInput struct {
	Slug     string
	Disabled *bool
	Admin    *bool
}

// DeleteSessionsForUserSlug finds all existing sessions for a given user, then removes them, effectively
// logging the user out of the service.
func DeleteSessionsForUserSlug(ctx context.Context, db *database.Connection, userSlug string) error {
	if !middleware.IsAdmin(ctx) {
		return backend.UnauthorizedReadErr(fmt.Errorf("Requesting user is not an admin"))
	}

	userID, err := userSlugToUserID(db, userSlug)
	if err != nil {
		return backend.DatabaseErr(err)
	}

	return deleteSessionsForUserID(db, userID)
}

// SetUserFlags updates flags for the indicated user, namely: admin and disabled.
// Then removes all sessions for that user (logging them out)
//
// NOTE: The flag is to _disable_ the user, which prevents access. To enable a user, set Disabled=false
func SetUserFlags(ctx context.Context, db *database.Connection, i SetUserFlagsInput) error {
	if !middleware.IsAdmin(ctx) {
		return backend.UnauthorizedReadErr(fmt.Errorf("Requesting user is not an admin"))
	}

	targetUser, err := db.RetrieveUserWithAuthDataBySlug(i.Slug)
	if err != nil {
		return backend.DatabaseErr(err)
	}
	err = validateAdminCanModifyFlag(ctx, targetUser, i)
	if err != nil {
		return backend.BadInputErr(err, err.Error())
	}

	valuesToUpdate := map[string]interface{}{}

	if i.Disabled != nil {
		valuesToUpdate["disabled"] = *i.Disabled
	}
	if i.Admin != nil {
		valuesToUpdate["admin"] = *i.Admin
	}

	if len(valuesToUpdate) > 0 {
		err := db.Update(sq.Update("users").SetMap(valuesToUpdate).Where(sq.Eq{"slug": i.Slug}))
		if err != nil {
			return backend.DatabaseErr(err)
		}
		return deleteSessionsForUserID(db, targetUser.ID)
	}
	return nil
}

// Note: this should only ever be done behind an IsAdmin check
func deleteSessionsForUserID(db *database.Connection, userID int64) error {

	if err := db.Exec("DELETE FROM sessions WHERE user_id = ?", userID); err != nil {
		return backend.DatabaseErr(err)
	}
	return nil
}

// validateAdminCanModifyFlag does some checks to validate the logic/sanity of the request.
// Checks roughly include logic verifying that a user isn't elevating/demoting their own status,
// users aren't being given status that doesn't make sense (specifically: headless users cannot be admins)
func validateAdminCanModifyFlag(ctx context.Context, targetUser models.UserWithAuthData, flagsToUpdate SetUserFlagsInput) error {
	targetUserIsSelf := targetUser.ID == middleware.UserID(ctx)
	targetUserIsHeadless := targetUser.Headless

	if flagsToUpdate.Admin != nil {
		// Note on valueUpdated: the frontend will supply all values, so an extra check here is done
		// to verify the intent to change the value rather than enforcing that no value is sent, without
		// requiring that the frontend explicitly omit values
		valueUpdated := targetUser.Admin != *flagsToUpdate.Admin
		if targetUserIsSelf && valueUpdated {
			return errors.New("Admins cannot alter their own admin status")
		}
		if targetUserIsHeadless && *flagsToUpdate.Admin == true {
			return errors.New("Headless users cannot be granted admin status")
		}
	}

	if flagsToUpdate.Disabled != nil {
		valueUpdated := targetUser.Disabled != *flagsToUpdate.Disabled
		if targetUserIsSelf && valueUpdated {
			return errors.New("Admins cannot disable themselves")
		}
	}

	return nil
}
