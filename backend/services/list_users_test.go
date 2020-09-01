// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestListUsers(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	testListUsersCase(t, db, "harry potter", true, []models.User{UserHarry})
	testListUsersCase(t, db, "granger", true, []models.User{UserHermione})
	testListUsersCase(t, db, "al", true, []models.User{UserAlastor, UserDumbledore, UserDraco, UserLucius, UserMinerva})
	testListUsersCase(t, db, "dra mal", true, []models.User{UserDraco})
	testListUsersCase(t, db, "", true, []models.User{})
	testListUsersCase(t, db, "  ", true, []models.User{})
	testListUsersCase(t, db, "%", true, []models.User{})
	testListUsersCase(t, db, "*", true, []models.User{})
	testListUsersCase(t, db, "___", true, []models.User{})

	// test for deleted user filtering
	testListUsersCase(t, db, UserTomRiddle.LastName, true, []models.User{UserTomRiddle})
	testListUsersCase(t, db, UserTomRiddle.LastName, false, []models.User{})
}

func testListUsersCase(t *testing.T, db *database.Connection, query string, includeDeleted bool, expectedUsers []models.User) {
	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})

	users, err := services.ListUsers(ctx, db, services.ListUsersInput{Query: query, IncludeDeleted: includeDeleted})
	require.NoError(t, err)

	require.Equal(t, len(expectedUsers), len(users), "Expected %d users for query '%s' but got %d", len(expectedUsers), query, len(users))

	for i := range expectedUsers {
		require.Equal(t, expectedUsers[i].Slug, users[i].Slug)
		require.Equal(t, expectedUsers[i].FirstName, users[i].FirstName)
		require.Equal(t, expectedUsers[i].LastName, users[i].LastName)
	}
}
