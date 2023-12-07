package services_test

import (
	"context"
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"
)

func TestCreateAPIKey(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserHermione
		targetUser := UserNeville
		adminUser := UserDumbledore
		ctx := contextForUser(normalUser, db)

		// Verify self actions
		verifyCreateAPIKey(t, false, ctx, db, normalUser.ID, "")
		verifyCreateAPIKey(t, false, ctx, db, normalUser.ID, normalUser.Slug)

		// verify other-based actions (non-admin)
		verifyCreateAPIKey(t, true, ctx, db, targetUser.ID, targetUser.Slug)

		// verify other-based actions (admin)
		ctx = contextForUser(adminUser, db)
		verifyCreateAPIKey(t, false, ctx, db, targetUser.ID, targetUser.Slug)
	})
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
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		targetUser := UserHarry
		adminUser := UserDumbledore
		ctx := contextForUser(normalUser, db)

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
		ctx = contextForUser(adminUser, db)
		verifyDeleteAPIKey(t, false, ctx, db, targetUser.ID, services.DeleteAPIKeyInput{
			UserSlug:  targetUser.Slug,
			AccessKey: APIKeyHarry1.AccessKey,
		})
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
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		targetUser := UserHarry
		adminUser := UserDumbledore
		ctx := contextForUser(normalUser, db)

		// verify read-self
		verifyListAPIKeys(t, false, ctx, db, "", APIKeyRon1, APIKeyRon2)
		// verify read-self (alt)
		verifyListAPIKeys(t, false, ctx, db, normalUser.Slug, APIKeyRon1, APIKeyRon2)

		// verify read-other (non-admin)
		verifyListAPIKeys(t, true, ctx, db, targetUser.Slug)

		// verify read-other (admin)
		ctx = contextForUser(adminUser, db)
		verifyListAPIKeys(t, false, ctx, db, targetUser.Slug, APIKeyHarry1, APIKeyHarry2)
	})
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
