// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	localConsts "github.com/theparanoids/ashirt-server/backend/authschemes/localauth/constants"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteAuthScheme(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	targetUser := UserHarry
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)
	schemeName := localConsts.Code
	recoveryScheme := recoveryConsts.Code

	// add some recovery data
	_, err := db.Insert("auth_scheme_data", map[string]interface{}{
		"auth_scheme": recoveryScheme,
		"user_key":    normalUser.FirstName,
		"user_id":     normalUser.ID,
	})
	require.NoError(t, err)
	_, err = db.Insert("auth_scheme_data", map[string]interface{}{
		"auth_scheme": recoveryScheme,
		"user_key":    targetUser.FirstName,
		"user_id":     targetUser.ID,
	})
	require.NoError(t, err)

	// verify cannot delete recovery
	verifyDeletedScheme(t, true, ctx, db, normalUser.ID, services.DeleteAuthSchemeInput{
		SchemeName: recoveryConsts.Code,
	})

	// verify delete auth for self
	verifyDeletedScheme(t, false, ctx, db, normalUser.ID, services.DeleteAuthSchemeInput{
		SchemeName: schemeName,
	})

	// add back in deleted auth
	_, err = db.Insert("auth_scheme_data", map[string]interface{}{
		"auth_scheme": schemeName,
		"user_key":    normalUser.FirstName,
		"user_id":     normalUser.ID,
	})
	require.NoError(t, err)

	// verify delete auth for self (alt)
	verifyDeletedScheme(t, false, ctx, db, normalUser.ID, services.DeleteAuthSchemeInput{
		UserSlug:   normalUser.Slug,
		SchemeName: schemeName,
	})

	// verify delete auth for other (non-admin)
	verifyDeletedScheme(t, true, ctx, db, targetUser.ID, services.DeleteAuthSchemeInput{
		UserSlug:   targetUser.Slug,
		SchemeName: schemeName,
	})

	// verify delete auth for other (admin)
	ctx = simpleFullContext(adminUser)
	verifyDeletedScheme(t, false, ctx, db, targetUser.ID, services.DeleteAuthSchemeInput{
		UserSlug:   targetUser.Slug,
		SchemeName: schemeName,
	})

	// verify cannot delete recovery (admin)
	verifyDeletedScheme(t, true, ctx, db, normalUser.ID, services.DeleteAuthSchemeInput{
		UserSlug:   targetUser.Slug,
		SchemeName: recoveryConsts.Code,
	})

}

func verifyDeletedScheme(t *testing.T, expectError bool, ctx context.Context, db *database.Connection, userID int64, input services.DeleteAuthSchemeInput) {
	originalSchemes := getAuthsForUser(t, db, userID)
	err := services.DeleteAuthScheme(ctx, db, input)
	if expectError {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)

	userSchemes := getAuthsForUser(t, db, userID)
	require.Equal(t, len(originalSchemes)-len(userSchemes), 1)
	for _, scheme := range userSchemes {
		require.NotEqual(t, scheme.AuthScheme, input.SchemeName)
	}
}
