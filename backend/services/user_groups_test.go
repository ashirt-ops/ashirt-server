// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/services"
)

type userGroupValidator func(*testing.T, UserOpPermJoinUser, *dtos.UserOperationRole)

// TODO TN
// ADD SEEDING and make specific tests instead of one big one
func TestCreateAndDeleteUserGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		i := services.ModifyUserGroupInput{
			Slug: "testGroup",
			UserSlugs: []string{
				UserRon.Slug,
				UserAlastor.Slug,
				UserHagrid.Slug,
			},
		}

		createUserGroupOutput, err := services.CreateUserGroup(db, i)
		require.NoError(t, err)
		require.Equal(t, createUserGroupOutput.RealSlug, i.Slug)
		userIDs, err := services.GetUserIDsFromGroup(db, createUserGroupOutput.UserGroupID)
		require.NoError(t, err)

		require.Equal(t, 3, len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserRon.ID, UserAlastor.ID, UserHagrid.ID}, userID)
		}

		createUserGroupOutput, err = services.CreateUserGroup(db, i)
		require.NoError(t, err)
		// Since a user group with that name already exists, a new slug should be created
		require.NotEqual(t, i.Slug, createUserGroupOutput.RealSlug)
		require.Contains(t, createUserGroupOutput.RealSlug, i.Slug)
		newUserIDs, _ := services.GetUserIDsFromGroup(db, createUserGroupOutput.UserGroupID)

		require.Equal(t, 3, len(newUserIDs))
		for _, userID := range newUserIDs {
			require.Contains(t, []int64{UserRon.ID, UserAlastor.ID, UserHagrid.ID}, userID)
		}

		err = services.DeleteUserGroup(db, createUserGroupOutput.RealSlug)
		require.NoError(t, err)
		userIDs, err = services.GetUserIDsFromGroup(db, createUserGroupOutput.UserGroupID)
		require.NoError(t, err)

		require.Equal(t, 0, len(userIDs))
	})
}

// func validateUserGroup(t *testing.T, expected UserOpPermJoinUser, actual *dtos.UserOperationRole) {
// 	require.Equal(t, expected.Slug, actual.User.Slug)
// 	require.Equal(t, expected.FirstName, actual.User.FirstName)
// 	require.Equal(t, expected.LastName, actual.User.LastName)
// 	require.Equal(t, expected.Role, actual.Role)
// }
