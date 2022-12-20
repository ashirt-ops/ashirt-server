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
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

type userGroupValidator func(*testing.T, UserOpPermJoinUser, *dtos.UserOperationRole)

func getUserIDsFromGroup(db *database.Connection, groupSlug string) ([]int64, error) {
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

func TestAddUsersToGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		gryffindorUserGroup := UserGroupGryffindor

		usersToAdd := []string{
			UserAlastor.Slug,
			UserHagrid.Slug,
		}

		err := services.AddUsersToGroup(db, usersToAdd, gryffindorUserGroup.ID)
		require.NoError(t, err)

		userIDs, err := getUserIDsFromGroup(db, gryffindorUserGroup.Slug)
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

		userIDs, err := getUserIDsFromGroup(db, gryffindorUserGroup.Slug)
		require.NoError(t, err)
		require.Equal(t, 2, len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserHarry.ID, UserGinny.ID}, userID)
		}
	})
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

		nonAdminUser := UserRon
		ctx := contextForUser(nonAdminUser, db)

		_, err := services.CreateUserGroup(ctx, db, i)
		// verify that non-admin user cannot create user groups
		require.Error(t, err)

		adminUser := UserDumbledore
		ctx = contextForUser(adminUser, db)
		_, err = services.CreateUserGroup(ctx, db, i)
		require.NoError(t, err)

		userIDs, err := getUserIDsFromGroup(db, slug)
		require.NoError(t, err)
		require.Equal(t, len(userSlugs), len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserRon.ID, UserAlastor.ID, UserHagrid.ID}, userID)
		}
		_, err = services.CreateUserGroup(ctx, db, i)
		assert.ErrorContains(t, err, "Unable to create user group. User group slug already exists")
	})
}

// TODO TN figure out why this test is so slow?
// probably the same reason why editing boht users and name at once doesn't work!!
func TestModifyUserGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		nonAdminUser := UserRon
		ctx := contextForUser(nonAdminUser, db)

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
		// verify that non-admin user cannot modify a user group
		require.Error(t, err)

		adminUser := UserDumbledore
		ctx = contextForUser(adminUser, db)

		_, err = services.ModifyUserGroup(ctx, db, i)
		require.NoError(t, err)

		userIDs, err := getUserIDsFromGroup(db, gryffindorUserGroup.Slug)
		require.NoError(t, err)
		require.Equal(t, 4, len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserHarry.ID, UserAlastor.ID, UserHagrid.ID, UserGinny.ID}, userID)
		}
	})
}

func TestDeleteUserGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		nonAdminUser := UserRon
		ctx := contextForUser(nonAdminUser, db)
		userGroup := UserGroupGryffindor

		err := services.DeleteUserGroup(ctx, db, userGroup.Slug)
		// verify that non-admin user cannot delete a user group
		require.Error(t, err)

		adminUser := UserDumbledore
		ctx = contextForUser(adminUser, db)

		err = services.DeleteUserGroup(ctx, db, userGroup.Slug)
		require.NoError(t, err)

		userIDs, err := getUserIDsFromGroup(db, userGroup.Slug)
		require.NoError(t, err)
		// 4 users in UserGroupGryffindor
		require.Equal(t, 4, len(userIDs))
	})
}

func TestListUserGroupsForAdmin(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		nonAdminUser := UserRon
		ctx := contextForUser(nonAdminUser, db)

		i := services.ListUserGroupsForAdminInput{
			Pagination: services.Pagination{
				TotalCount: 4,
				PageSize:   10,
				Page:       1,
			},
			IncludeDeleted: false,
		}

		result, err := services.ListUserGroupsForAdmin(ctx, db, i)
		// verify that non-admin user cannot list user groups
		require.Error(t, err)

		adminUser := UserDumbledore
		ctx = contextForUser(adminUser, db)

		result, err = services.ListUserGroupsForAdmin(ctx, db, i)
		var usergroups = result.Content.([]dtos.UserGroupAdminView)
		require.Equal(t, int64(1), result.PageNumber)
		require.Equal(t, int64(4), result.PageSize)
		require.Equal(t, int64(4), result.TotalCount)
		require.Equal(t, 4, len(usergroups))
		require.NoError(t, err)
	})
}

func TestListUserGroupsForOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpSorcerersStone
		allUserGroupOpRoles := getUserGroupsWithRoleForOperationByOperationID(t, db, masterOp.ID)
		require.NotEqual(t, len(allUserGroupOpRoles), 0, "Some user groups should be attached to this operation")

		input := services.ListUserGroupsForOperationInput{
			OperationSlug: masterOp.Slug,
		}

		content, err := services.ListUserGroupsForOperation(ctx, db, input)
		// Ron is not an operation admin, so he should not be able to list user groups
		require.Error(t, err)

		ctx = contextForUser(UserHarry, db)
		content, err = services.ListUserGroupsForOperation(ctx, db, input)
		require.NoError(t, err)

		require.Equal(t, len(content), len(allUserGroupOpRoles))
		validateUserGroupSets(t, content, allUserGroupOpRoles)
	})
}

func TestListUserGroups(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		testListUserGroupsCase(t, db, "gryf", true, []models.UserGroup{UserGroupGryffindor})
		testListUserGroupsCase(t, db, "ff", true, []models.UserGroup{UserGroupGryffindor, UserGroupHufflepuff})
		testListUserGroupsCase(t, db, "l", true, []models.UserGroup{UserGroupHufflepuff, UserGroupRavenclaw, UserGroupSlytherin})
		testListUserGroupsCase(t, db, "", true, []models.UserGroup{})
		testListUserGroupsCase(t, db, "  ", true, []models.UserGroup{})
		testListUserGroupsCase(t, db, "%", true, []models.UserGroup{})
		testListUserGroupsCase(t, db, "*", true, []models.UserGroup{})
		testListUserGroupsCase(t, db, "___", true, []models.UserGroup{})

		// test for deleted user filtering
		testListUserGroupsCase(t, db, UserGroupOtherHouse.Name, true, []models.UserGroup{UserGroupOtherHouse})
		testListUserGroupsCase(t, db, UserTomRiddle.LastName, false, []models.UserGroup{})
	})
}

func testListUserGroupsCase(t *testing.T, db *database.Connection, query string, includeDeleted bool, expectedUserGroups []models.UserGroup) {
	ctx := contextForUser(UserDumbledore, db)

	userGroups, err := services.ListUserGroups(ctx, db, services.ListUserGroupsInput{Query: query, IncludeDeleted: includeDeleted})
	require.NoError(t, err)

	require.Equal(t, len(expectedUserGroups), len(userGroups), "Expected %d users for query '%s' but got %d", len(expectedUserGroups), query, len(userGroups))

	for i := range expectedUserGroups {
		require.Equal(t, expectedUserGroups[i].Slug, userGroups[i].Slug)
		require.Equal(t, expectedUserGroups[i].Name, userGroups[i].Name)
	}
}

func validateUserGroupSets(t *testing.T, dtoSet []*dtos.UserGroupOperationRole, dbSet []UserGroupOpPermJoinUser) {
	var expected *UserGroupOpPermJoinUser = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.Slug == dtoItem.UserGroup.Slug {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		require.Equal(t, expected.Slug, dtoItem.UserGroup.Slug)
		require.Equal(t, expected.Name, dtoItem.UserGroup.Name)
		require.Equal(t, expected.Role, dtoItem.Role)
	}
}
