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
