// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"

	sq "github.com/Masterminds/squirrel"
)

func TestDeleteSessionsForUserSlug(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	targetedUser := UserDraco
	alsoPresentUser := UserHarry

	// populate some sessions
	sessionsToAdd := []models.Session{
		models.Session{UserID: targetedUser.ID, SessionData: []byte("a")},
		models.Session{UserID: alsoPresentUser.ID, SessionData: []byte("b")},
		models.Session{UserID: targetedUser.ID, SessionData: []byte("c")},
		models.Session{UserID: alsoPresentUser.ID, SessionData: []byte("d")},
	}
	err := db.BatchInsert("sessions", len(sessionsToAdd), func(i int) map[string]interface{} {
		return map[string]interface{}{
			"user_id":      sessionsToAdd[i].UserID,
			"session_data": sessionsToAdd[i].SessionData,
		}
	})

	require.NoError(t, err)

	// verify sessions exist
	var targetedUserSessions []models.Session
	err = db.Select(&targetedUserSessions, sq.Select("*").From("sessions").Where(sq.Eq{"user_id": targetedUser.ID}))
	require.NoError(t, err)
	require.True(t, len(targetedUserSessions) > 0)

	// verify non-admin cannot delete session data
	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	err = services.DeleteSessionsForUserSlug(ctx, db, targetedUser.Slug)
	require.Error(t, err)

	// verify admin can delete session data
	ctx = fullContextAsAdmin(UserDumbledore.ID, &policy.FullAccess{})
	err = services.DeleteSessionsForUserSlug(ctx, db, targetedUser.Slug)
	require.NoError(t, err)

	targetedUserSessions = []models.Session{}
	err = db.Select(&targetedUserSessions, sq.Select("*").From("sessions").Where(sq.Eq{"user_id": targetedUser.ID}))
	require.NoError(t, err)
	require.True(t, len(targetedUserSessions) == 0)
}

func TestSetUserFlags(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	targetUser := UserHarry
	adminUser := UserDumbledore
	admin := true
	disabled := true
	input := services.SetUserFlagsInput{
		Slug:     targetUser.Slug,
		Admin:    &admin,
		Disabled: &disabled,
	}

	// verify access restricted for non-admins
	ctx := fullContext(UserDraco.ID, &policy.FullAccess{}) // Note: not an admin
	err := services.SetUserFlags(ctx, db, input)
	require.Error(t, err)

	// As an admin
	ctx = fullContextAsAdmin(adminUser.ID, &policy.FullAccess{})

	// verify users can't disable themselves
	sameUserInput := services.SetUserFlagsInput{
		Slug:     adminUser.Slug,
		Admin:    &admin,    // true at this point (no change)
		Disabled: &disabled, // true at this point
	}
	err = services.SetUserFlags(ctx, db, sameUserInput)
	require.Error(t, err)

	// verify users can't demote themselves
	disabled = false
	admin = false
	err = services.SetUserFlags(ctx, db, sameUserInput)
	require.Error(t, err)

	// reset for next tests
	disabled = true
	admin = true

	// try setting and then unsetting admin/disabled
	for i := 0; i < 2; i++ {
		err = services.SetUserFlags(ctx, db, input)
		require.NoError(t, err)

		dbProfile := getUserProfile(t, db, targetUser.ID)

		require.Equal(t, admin, dbProfile.Admin)
		require.Equal(t, disabled, dbProfile.Disabled)

		// second test: Make sure setting to false also works
		admin = !admin
		disabled = !disabled
	}

	// verify headless users cannot be admins
	admin = true
	err = services.SetUserFlags(ctx, db, services.SetUserFlagsInput{
		Slug:  UserHeadlessNick.Slug,
		Admin: &admin,
	})
	require.Error(t, err)
}
