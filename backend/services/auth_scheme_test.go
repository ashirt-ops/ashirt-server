// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	localConsts "github.com/theparanoids/ashirt-server/backend/authschemes/localauth/constants"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteAuthScheme(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		targetUser := UserHarry
		ctx := contextForUser(normalUser, db)
		schemeName := localConsts.Code
		recoveryScheme := recoveryConsts.Code

		// seed recovery data
		users := []models.User{normalUser, targetUser}
		err := db.BatchInsert("auth_scheme_data", len(users), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"auth_scheme": recoveryScheme,
				"auth_type":   recoveryScheme,
				"user_key":    users[i].FirstName,
				"user_id":     users[i].ID,
			}
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
			"auth_type":   schemeName,
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
		ctx = contextForUser(UserDumbledore, db)
		verifyDeletedScheme(t, false, ctx, db, targetUser.ID, services.DeleteAuthSchemeInput{
			UserSlug:   targetUser.Slug,
			SchemeName: schemeName,
		})

		// verify cannot delete recovery (admin)
		verifyDeletedScheme(t, true, ctx, db, normalUser.ID, services.DeleteAuthSchemeInput{
			UserSlug:   targetUser.Slug,
			SchemeName: recoveryConsts.Code,
		})
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

func TestDeleteAuthSchemeUsers(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		adminUser := UserDumbledore
		schemeName := localConsts.Code

		baseUsers := getUsersForAuth(t, db, schemeName)
		require.Greater(t, len(baseUsers), 0)

		// verify non-admins have no access
		ctx := contextForUser(normalUser, db)
		err := services.DeleteAuthSchemeUsers(ctx, db, schemeName)
		require.Error(t, err)

		// verify admins have access + effect works
		ctx = contextForUser(adminUser, db)
		err = services.DeleteAuthSchemeUsers(ctx, db, schemeName)
		require.NoError(t, err)

		updatedUsers := getUsersForAuth(t, db, schemeName)
		require.Equal(t, 0, len(updatedUsers))

		// verify admins cannot delete recovery
		err = services.DeleteAuthSchemeUsers(ctx, db, recoveryConsts.Code)
		require.Error(t, err)
	})
}

var patronusAuthScheme = dtos.SupportedAuthScheme{SchemeName: "Patronus Charm", SchemeCode: "patronus", SchemeType: "magical"}
var localAuthScheme = dtos.SupportedAuthScheme{SchemeName: localConsts.FriendlyName, SchemeCode: localConsts.Code, SchemeType: localConsts.Code}
var darkMarkScheme = dtos.SupportedAuthScheme{SchemeName: "Death Eaters", SchemeCode: "dark mark", SchemeType: "magical"}

func TestListAuthDetailsKeys(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		adminUser := UserDumbledore

		supportedSchemes := []dtos.SupportedAuthScheme{
			localAuthScheme,
			patronusAuthScheme,
		}

		// verify non-admins cannot access this service
		ctx := contextForUser(normalUser, db)
		_, err := services.ListAuthDetails(ctx, db, &supportedSchemes)
		require.Error(t, err)

		// verify list for admins
		ctx = contextForUser(adminUser, db)
		results, err := services.ListAuthDetails(ctx, db, &supportedSchemes)
		require.NoError(t, err)
		require.Equal(t, 2, len(results))

		var schemePatronus *dtos.DetailedAuthenticationInfo
		var schemeLocal *dtos.DetailedAuthenticationInfo
		for _, scheme := range results {
			if scheme.AuthSchemeCode == localAuthScheme.SchemeCode {
				schemeLocal = scheme
			}
			if scheme.AuthSchemeCode == patronusAuthScheme.SchemeCode {
				schemePatronus = scheme
			}
		}
		// verify expected results
		require.NotNil(t, schemePatronus)
		require.NotNil(t, schemeLocal)

		require.Equal(t, int64(0), schemePatronus.UserCount)
		require.Equal(t, int64(0), schemePatronus.UniqueUserCount)
		require.Equal(t, 0, len(schemePatronus.Labels))

		realUserCount := int64(len(getRealUsers(t, db)))
		require.Equal(t, realUserCount, schemeLocal.UserCount)
		require.Equal(t, realUserCount, schemeLocal.UniqueUserCount)
		require.Equal(t, 0, len(schemeLocal.Labels))

		// Add in an unsupported scheme
		_, err = db.Insert("auth_scheme_data", map[string]interface{}{
			"auth_scheme": darkMarkScheme.SchemeCode,
			"user_key":    "Half-Blood Prince",
			"user_id":     UserSnape.ID,
			"auth_type":   darkMarkScheme.SchemeType,
		})

		require.NoError(t, err)
		results, err = services.ListAuthDetails(ctx, db, &supportedSchemes)
		require.NoError(t, err)
		require.Equal(t, 3, len(results))
		var schemeDarkMark *dtos.DetailedAuthenticationInfo

		for _, scheme := range results {
			if scheme.AuthSchemeCode == darkMarkScheme.SchemeCode {
				schemeDarkMark = scheme
			}
			if scheme.AuthSchemeCode == localAuthScheme.SchemeCode {
				schemeLocal = scheme
			}
		}
		require.NotNil(t, schemeDarkMark)

		// verify count whent down for local
		require.Equal(t, realUserCount, schemeLocal.UserCount)
		require.Equal(t, realUserCount-1, schemeLocal.UniqueUserCount)
		require.Equal(t, 0, len(schemeLocal.Labels))

		require.Equal(t, int64(1), schemeDarkMark.UserCount)
		require.Equal(t, int64(0), schemeDarkMark.UniqueUserCount)
		require.Equal(t, 1, len(schemeDarkMark.Labels))
		require.Equal(t, "Unsupported", schemeDarkMark.Labels[0])
	})
}
