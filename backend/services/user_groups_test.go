// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/services"
)

type userGroupValidator func(*testing.T, UserOpPermJoinUser, *dtos.UserOperationRole)

func GetUserIDsFromGroup(db *database.Connection, groupSlug string) ([]int64, error) {
	var userGroupId int64
	err := db.Get(&userGroupId, sq.Select("id").
		From("user_groups").
		Where(sq.Eq{
			"slug": groupSlug,
		}))
	if err != nil {
		s := fmt.Sprintf("Cannot get user group id for group %q", groupSlug)
		return nil, backend.WrapError(s, backend.DatabaseErr(err))
	}

	var userGroupMap []int64
	err = db.Select(&userGroupMap, sq.Select("user_id").
		From("group_user_map").
		Where(sq.Eq{
			"group_id": userGroupId,
		}))
	if err != nil {
		s := fmt.Sprintf("Cannot get user group map for group %q", userGroupId)
		return userGroupMap, backend.WrapError(s, backend.DatabaseErr(err))
	}
	return userGroupMap, nil
}

func TestCreateUserGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		slug := "testGroup"
		userSlugs := []string{
			UserRon.Slug,
			UserAlastor.Slug,
			UserHagrid.Slug,
		}
		i := services.CreateUserGroupInput{
			Name:      slug,
			Slug:      slug,
			UserSlugs: userSlugs,
		}

		adminUser := UserDumbledore
		ctx := contextForUser(adminUser, db)
		_, err := services.CreateUserGroup(ctx, db, i)
		require.NoError(t, err)

		userIDs, err := GetUserIDsFromGroup(db, slug)
		require.NoError(t, err)
		require.Equal(t, len(userSlugs), len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserRon.ID, UserAlastor.ID, UserHagrid.ID}, userID)
		}
		_, err = services.CreateUserGroup(ctx, db, i)
		assert.ErrorContains(t, err, "Unable to create user group. User group slug already exists")
	})
}

func TestDeleteUserGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		adminUser := UserDumbledore
		ctx := contextForUser(adminUser, db)
		userGroup := UserGroupGryffindor

		err := services.DeleteUserGroup(ctx, db, userGroup.Slug)
		require.NoError(t, err)

		userIDs, err := GetUserIDsFromGroup(db, userGroup.Slug)
		require.NoError(t, err)
		// 4 users in UserGroupGryffindor
		require.Equal(t, 4, len(userIDs))
	})
}

// TODO TN figure out why this test is so slow?
func TestModifyUserGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		adminUser := UserDumbledore
		ctx := contextForUser(adminUser, db)
		gryffindorUserGroup := UserGroupGryffindor

		newName := "Glyssintor"
		usersToAdd := []string{
			UserAlastor.Slug,
			UserHagrid.Slug,
		}
		usersToRemove := []string{
			UserRon.Slug,
			UserHermione.Slug,
		}
		i := services.ModifyUserGroupInput{
			Name:          newName,
			Slug:          gryffindorUserGroup.Slug,
			UsersToAdd:    usersToAdd,
			UsersToRemove: usersToRemove,
		}
		// TODO TN check that name actually changed by grabbing record

		_, err := services.ModifyUserGroup(ctx, db, i)
		require.NoError(t, err)

		// userIDs, err := GetUserIDsFromGroup(db, gryffindorUserGroup.Slug)
		// require.NoError(t, err)
		// // TODO TN figure out why this is incorrect?
		// require.Equal(t, 4, len(userIDs))
		// for _, userID := range userIDs {
		// 	require.Contains(t, []int64{UserHarry.ID, UserAlastor.ID, UserHagrid.ID, UserGinny.ID}, userID)
		// }
	})
}

func TestListUserGroups(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		adminUser := UserDumbledore
		ctx := contextForUser(adminUser, db)

		i := services.ListUserGroupsForAdminInput{
			Pagination: services.Pagination{
				TotalCount: 4,
				PageSize:   10,
				Page:       1,
			},
			IncludeDeleted: false,
		}

		result, err := services.ListUserGroupsForAdmin(ctx, db, i)
		var usergroups = result.Content.([]dtos.UserGroupAdminView)
		require.Equal(t, int64(1), result.PageNumber)
		require.Equal(t, int64(5), result.PageSize)
		require.Equal(t, int64(5), result.TotalCount)
		require.Equal(t, 5, len(usergroups))
		require.NoError(t, err)
	})
}

// write a test to test AddUsersToGroup and RemoveUsersFromGroup

func TestAddUsersToGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		gryffindorUserGroup := UserGroupGryffindor

		usersToAdd := []string{
			UserAlastor.Slug,
			UserHagrid.Slug,
		}

		err := services.AddUsersToGroup(db, usersToAdd, gryffindorUserGroup.ID)
		require.NoError(t, err)

		userIDs, err := GetUserIDsFromGroup(db, gryffindorUserGroup.Slug)
		require.NoError(t, err)
		require.Equal(t, 6, len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserHarry.ID, UserRon.ID, UserHermione.ID, UserAlastor.ID, UserHagrid.ID, UserGinny.ID}, userID)
		}
	})
}

func TestRemoveUsersFromGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		gryffindorUserGroup := UserGroupGryffindor

		usersToRemove := []string{
			UserRon.Slug,
			UserHermione.Slug,
		}

		err := services.RemoveUsersFromGroup(db, usersToRemove, gryffindorUserGroup.ID)
		require.NoError(t, err)

		userIDs, err := GetUserIDsFromGroup(db, gryffindorUserGroup.Slug)
		require.NoError(t, err)
		require.Equal(t, 2, len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserHarry.ID, UserGinny.ID}, userID)
		}
	})
}
