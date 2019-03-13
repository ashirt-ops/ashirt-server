// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"testing"

	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

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
