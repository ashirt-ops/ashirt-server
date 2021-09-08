// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	localConsts "github.com/theparanoids/ashirt-server/backend/authschemes/localauth/constants"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/services"
)

var patronusAuthScheme = dtos.SupportedAuthScheme{SchemeName: "Patronus Charm", SchemeCode: "patronus", SchemeType: "magical"}
var localAuthScheme = dtos.SupportedAuthScheme{SchemeName: localConsts.FriendlyName, SchemeCode: localConsts.Code, SchemeType: localConsts.Code}
var darkMarkScheme = dtos.SupportedAuthScheme{SchemeName: "Death Eaters", SchemeCode: "dark mark", SchemeType: "magical"}

func TestListAuthDetailsKeys(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)

	supportedSchemes := []dtos.SupportedAuthScheme{
		localAuthScheme,
		patronusAuthScheme,
	}

	// verify non-admins cannot access this service
	_, err := services.ListAuthDetails(ctx, db, &supportedSchemes)
	require.Error(t, err)

	// verify list for admins
	ctx = simpleFullContext(adminUser)
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
}
