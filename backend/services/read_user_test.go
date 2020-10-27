// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestReadUser(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	targetUser := UserHarry
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)

	supportedAuthSchemes := []dtos.SupportedAuthScheme{
		dtos.SupportedAuthScheme{SchemeName: "Local", SchemeCode: "local"},
	}

	// verify read-self
	retrievedUser, err := services.ReadUser(ctx, db, "", &supportedAuthSchemes)
	require.NoError(t, err)
	verifyRetrievedUser(t, normalUser, retrievedUser, supportedAuthSchemes)

	// verify read-self alternative (userslug provided)
	retrievedUser, err = services.ReadUser(ctx, db, normalUser.Slug, &supportedAuthSchemes)
	require.NoError(t, err)
	verifyRetrievedUser(t, normalUser, retrievedUser, supportedAuthSchemes)

	// verify read-other (non-admin : should fail)
	_, err = services.ReadUser(ctx, db, targetUser.Slug, &supportedAuthSchemes)
	require.Error(t, err)

	// verify read-other (as admin)
	ctx = simpleFullContext(adminUser)
	retrievedUser, err = services.ReadUser(ctx, db, targetUser.Slug, &supportedAuthSchemes)
	require.NoError(t, err)
	verifyRetrievedUser(t, targetUser, retrievedUser, supportedAuthSchemes)

	// verify old/removed auth schemes are filtered out
	ctx = simpleFullContext(normalUser)
	supportedAuthSchemes = []dtos.SupportedAuthScheme{
		dtos.SupportedAuthScheme{SchemeName: "Petronus", SchemeCode: "petroni"},
	}
	retrievedUser, err = services.ReadUser(ctx, db, "", &supportedAuthSchemes)
	require.NoError(t, err)
	verifyRetrievedUser(t, normalUser, retrievedUser, []dtos.SupportedAuthScheme{})
}

func verifyRetrievedUser(t *testing.T, expectedUser models.User, retrievedUser *dtos.UserOwnView, expectedAuths []dtos.SupportedAuthScheme) {
	require.Equal(t, expectedUser.Slug, retrievedUser.Slug)
	require.Equal(t, expectedUser.FirstName, retrievedUser.FirstName)
	require.Equal(t, expectedUser.LastName, retrievedUser.LastName)
	require.Equal(t, expectedUser.Email, retrievedUser.Email)
	for _, expectedAuth := range expectedAuths {
		found := false

		for _, returnedAuth := range retrievedUser.Authentication {
			if expectedAuth.SchemeCode == returnedAuth.AuthSchemeCode {
				found = true
				break
			}
		}
		require.True(t, found)
	}
}
