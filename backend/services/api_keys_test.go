// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateAPIKey(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserHermione
	targetUser := UserNeville
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)

	// Verify self actions
	verifyCreateAPIKey(t, false, ctx, db, normalUser.ID, "")
	verifyCreateAPIKey(t, false, ctx, db, normalUser.ID, normalUser.Slug)

	// verify other-based actions (non-admin)
	verifyCreateAPIKey(t, true, ctx, db, targetUser.ID, targetUser.Slug)

	// verify other-based actions (admin)
	ctx = simpleFullContext(adminUser)
	verifyCreateAPIKey(t, false, ctx, db, targetUser.ID, targetUser.Slug)
}

func verifyCreateAPIKey(t *testing.T, expectError bool, ctx context.Context, db *database.Connection, userID int64, userSlug string) {
	originalKeys := getAPIKeysForUserID(t, db, userID)
	apiKey, apiErr := services.CreateAPIKey(ctx, db, userSlug)
	if expectError {
		require.Error(t, apiErr)
		return
	}
	require.NoError(t, apiErr)

	require.Len(t, apiKey.AccessKey, 24)
	require.Len(t, apiKey.SecretKey, 64)

	userAPIKeys := getAPIKeysForUserID(t, db, userID)
	require.Equal(t, len(userAPIKeys)-len(originalKeys), 1)

	lastKey := userAPIKeys[len(userAPIKeys)-1]
	require.Equal(t, apiKey.AccessKey, lastKey.AccessKey)
	require.Equal(t, apiKey.SecretKey, lastKey.SecretKey)
}

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

func TestListAPIKeys(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	targetUser := UserHarry
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)

	// verify read-self
	verifyListAPIKeys(t, false, ctx, db, "", APIKeyRon1, APIKeyRon2)
	// verify read-self (alt)
	verifyListAPIKeys(t, false, ctx, db, normalUser.Slug, APIKeyRon1, APIKeyRon2)

	// verify read-other (non-admin)
	verifyListAPIKeys(t, true, ctx, db, targetUser.Slug)

	// verify read-other (admin)
	ctx = simpleFullContext(adminUser)
	verifyListAPIKeys(t, false, ctx, db, targetUser.Slug, APIKeyHarry1, APIKeyHarry2)
}

func verifyListAPIKeys(t *testing.T, expectError bool, ctx context.Context, db *database.Connection, userSlug string, expectedAPIKeys ...models.APIKey) {
	apiKeys, err := services.ListAPIKeys(ctx, db, userSlug)
	if expectError {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)

	require.Len(t, apiKeys, len(expectedAPIKeys))
	require.Equal(t, apiKeys[0].AccessKey, expectedAPIKeys[0].AccessKey)
	require.Equal(t, apiKeys[1].AccessKey, expectedAPIKeys[1].AccessKey)
	require.Nil(t, apiKeys[0].SecretKey)
	require.Nil(t, apiKeys[1].SecretKey)
}

func TestRotateAPIKey(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	targetUser := UserRon
	nonAdminUser := UserNeville
	adminUser := UserDumbledore

	// verify user can rotate their own keys
	ctx := contextForUser(targetUser, db)
	verifyRotateAPIKeys(t, false, ctx, db, targetUser.Slug, false) // primary
	verifyRotateAPIKeys(t, false, ctx, db, targetUser.Slug, true)  // Alt

	// Verify that admins can rotate keys for a user
	ctx = contextForUser(adminUser, db)
	verifyRotateAPIKeys(t, true, ctx, db, targetUser.Slug, false) // admins cannot change an api key without specifying the slug (unless they're themselves)
	verifyRotateAPIKeys(t, false, ctx, db, targetUser.Slug, true) // Alt

	// Verify that others cannot rotate an api key not owned by them
	originalAPIKey, err := services.CreateAPIKey(ctx, db, targetUser.Slug)
	require.NoError(t, err)

	ctx = contextForUser(nonAdminUser, db)
	_, err = services.RotateAPIKey(ctx, db, services.RotateAPIKeyInput{AccessKey: originalAPIKey.AccessKey})
	require.Error(t, err)
	_, err = services.RotateAPIKey(ctx, db, services.RotateAPIKeyInput{AccessKey: originalAPIKey.AccessKey, UserSlug: targetUser.Slug})
	require.Error(t, err)
}

// verifyRotateAPIKeys only works when the context reflects an admin user, or a normal user targeting themself
func verifyRotateAPIKeys(t *testing.T, expectError bool, ctx context.Context, db *database.Connection, apiKeyOwnerSlug string, withSlug bool) {
	// create initial key to test with
	originalAPIKey, err := services.CreateAPIKey(ctx, db, apiKeyOwnerSlug)
	require.NoError(t, err)

	// get a list of the new keys
	originalKeys, err := services.ListAPIKeys(ctx, db, apiKeyOwnerSlug)
	require.NoError(t, err)

	// rotate the newest key
	input := services.RotateAPIKeyInput{
		AccessKey: originalAPIKey.AccessKey,
	}
	if withSlug {
		input.UserSlug = apiKeyOwnerSlug
	}
	newAPIKey, err := services.RotateAPIKey(ctx, db, input)
	if expectError {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)

	// verify that the old key doesn't exist, but the new key does
	updatedKeys, err := services.ListAPIKeys(ctx, db, apiKeyOwnerSlug)
	require.NoError(t, err)
	require.Equal(t, len(originalKeys), len(updatedKeys)) // check verifies that one was deleted, one was added -- net 0 change

	originalKeyWasDeleted := true
	newKeyWasFound := false

	for _, key := range updatedKeys {
		if key.AccessKey == newAPIKey.AccessKey {
			newKeyWasFound = true
		}
		if key.AccessKey == originalAPIKey.AccessKey {
			originalKeyWasDeleted = false
		}
	}
	require.True(t, newKeyWasFound)
	require.True(t, originalKeyWasDeleted)
}
