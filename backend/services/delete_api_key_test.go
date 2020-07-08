// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteAPIKey(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	targetUser := UserHarry
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)

	// verify delete api key for other user (as self)
	verifyDeleteAPIKey(t, true, ctx, db, normalUser.ID, services.DeleteAPIKeyInput{
		UserSlug:  normalUser.Slug,
		AccessKey: APIKeyHarry1.AccessKey,
	})

	// verify delete api key for other user (as self - alt)
	verifyDeleteAPIKey(t, true, ctx, db, normalUser.ID, services.DeleteAPIKeyInput{
		AccessKey: APIKeyHarry1.AccessKey,
	})

	// verify delete api key for self
	verifyDeleteAPIKey(t, false, ctx, db, normalUser.ID, services.DeleteAPIKeyInput{
		AccessKey: APIKeyRon1.AccessKey,
	})

	// verify delete api key for self (alt)
	verifyDeleteAPIKey(t, false, ctx, db, normalUser.ID, services.DeleteAPIKeyInput{
		UserSlug:  normalUser.Slug,
		AccessKey: APIKeyRon2.AccessKey,
	})

	// verify delete api key for other (non-admin)
	verifyDeleteAPIKey(t, true, ctx, db, targetUser.ID, services.DeleteAPIKeyInput{
		UserSlug:  targetUser.Slug,
		AccessKey: APIKeyHarry1.AccessKey,
	})

	// verify delete api key for other (admin)
	ctx = simpleFullContext(adminUser)
	verifyDeleteAPIKey(t, false, ctx, db, targetUser.ID, services.DeleteAPIKeyInput{
		UserSlug:  targetUser.Slug,
		AccessKey: APIKeyHarry1.AccessKey,
	})
}

func verifyDeleteAPIKey(t *testing.T, expectError bool, ctx context.Context, db *database.Connection, userID int64, input services.DeleteAPIKeyInput) {
	originalKeys := getAPIKeysForUserID(t, db, userID)
	require.Greater(t, len(originalKeys), 0, "This test only works for users that have at least 1 key")
	err := services.DeleteAPIKey(ctx, db, input)
	if expectError {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)

	updatedKeys := getAPIKeysForUserID(t, db, userID)
	require.Equal(t, len(originalKeys)-len(updatedKeys), 1)
	i := 0
	for _, k := range updatedKeys {
		if k.AccessKey == input.AccessKey {
			continue
		}
		require.Equal(t, k.AccessKey, updatedKeys[i].AccessKey)
		i++
	}
}
