// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestUpdateUserProfile(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	targetUser := UserHarry
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)

	// verify read-self
	verifyUserProfileUpdate(t, false, ctx, db, normalUser.ID, services.UpdateUserProfileInput{
		FirstName: "Stan",
		LastName:  "Shunpike",
		Email:     "sshunpike@hogwarts.edu",
	})

	// verify read-self (alternate)
	verifyUserProfileUpdate(t, false, ctx, db, normalUser.ID, services.UpdateUserProfileInput{
		UserSlug:  normalUser.Slug,
		FirstName: "Stan2",
		LastName:  "Shunpike2",
		Email:     "sshunpike2@hogwarts.edu",
	})

	// verify read-other (non-admin)
	verifyUserProfileUpdate(t, true, ctx, db, targetUser.ID, services.UpdateUserProfileInput{
		UserSlug:  targetUser.Slug,
		FirstName: "Stan3",
		LastName:  "Shunpike3",
		Email:     "sshunpike3@hogwarts.edu",
	})

	// verify read-other (admin)
	ctx = simpleFullContext(adminUser)
	verifyUserProfileUpdate(t, false, ctx, db, targetUser.ID, services.UpdateUserProfileInput{
		UserSlug:  targetUser.Slug,
		FirstName: "Stan4",
		LastName:  "Shunpike4",
		Email:     "sshunpike4@hogwarts.edu",
	})

}

func verifyUserProfileUpdate(t *testing.T, expectError bool, ctx context.Context, db *database.Connection, userID int64, updatedData services.UpdateUserProfileInput) {
	err := services.UpdateUserProfile(ctx, db, updatedData)
	if expectError {
		require.NotNil(t, err)
		return
	}

	require.NoError(t, err)

	newProfile := getUserProfile(t, db, userID)
	require.NoError(t, err)
	require.Equal(t, updatedData.FirstName, newProfile.FirstName)
	require.Equal(t, updatedData.LastName, newProfile.LastName)
	require.Equal(t, updatedData.Email, newProfile.Email)
}
