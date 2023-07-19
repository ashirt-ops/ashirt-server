// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		err := db.WithTx(context.Background(), func(tx *database.Transactable) {
			services.AddUsersToGroup(tx, usersToAdd, gryffindorUserGroup.ID)
		})

		require.NoError(t, err)

		userIDs, err := getUserIDsFromGroup(db, gryffindorUserGroup.Slug)
		require.NoError(t, err)
		require.Equal(t, 6, len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserHarry.ID, UserRon.ID, UserHermione.ID, UserAlastor.ID, UserHagrid.ID, UserGinny.ID}, userID)
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
		assert.Error(t, err)
	})
}

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

		_, err := services.ModifyUserGroup(ctx, db, i)
		// verify that non-admin user cannot modify a user group
		require.Error(t, err)

		adminUser := UserDumbledore
		ctx = contextForUser(adminUser, db)

		result, err := services.ModifyUserGroup(ctx, db, i)
		require.NoError(t, err)
		fullUserGroup := getUserGroupFromSlug(t, db, result.Slug)
		require.Equal(t, newName, fullUserGroup.Name)

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
			IncludeDeleted: false,
		}

		_, err := services.ListUserGroupsForAdmin(ctx, db, i)
		// verify that non-admin user cannot list user groups
		require.Error(t, err)

		adminUser := UserDumbledore
		ctx = contextForUser(adminUser, db)

		_, err = services.ListUserGroupsForAdmin(ctx, db, i)
		require.NoError(t, err)
	})
}

func TestGetSlugMap(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		i := services.ListUserGroupsForAdminInput{
			IncludeDeleted: true,
		}

		slugMap, err := services.GetSlugMap(db, i)
		require.NoError(t, err)
		require.Equal(t, 12, len(slugMap))
		for _, slugMapEntry := range slugMap {
			userName := slugMapEntry.UserSlug.String
			if userName != "" {
				require.Contains(t, []string{UserHarry.Slug, UserGinny.Slug, UserRon.Slug, UserHermione.Slug, UserCedric.Slug, UserCho.Slug, UserFleur.Slug, UserViktor.Slug, UserSnape.Slug, UserLucius.Slug, UserDraco.Slug}, userName)
			}
			if slugMapEntry.Deleted.Valid == true {
				require.Equal(t, UserGroupOtherHouse.Slug, slugMapEntry.GroupSlug)
			}
		}

		// test for non-deleted user groups
		i = services.ListUserGroupsForAdminInput{
			IncludeDeleted: false,
		}

		slugMap, err = services.GetSlugMap(db, i)
		require.NoError(t, err)
		require.Equal(t, 11, len(slugMap))
	})
}

func TestSortUsersInToGroups(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		slugMap := services.SlugMap{
			{
				UserSlug: sql.NullString{
					String: UserHarry.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupGryffindor.Slug,
				GroupName: UserGroupGryffindor.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserRon.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupGryffindor.Slug,
				GroupName: UserGroupGryffindor.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserGinny.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupGryffindor.Slug,
				GroupName: UserGroupGryffindor.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserHermione.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupGryffindor.Slug,
				GroupName: UserGroupGryffindor.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserCedric.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupHufflepuff.Slug,
				GroupName: UserGroupHufflepuff.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserFleur.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupHufflepuff.Slug,
				GroupName: UserGroupHufflepuff.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			// Includes groups without a user, we need to return those groups as well
			{
				UserSlug: sql.NullString{
					String: "",
					Valid:  false,
				},
				GroupSlug: UserGroupOtherHouse.Slug,
				GroupName: UserGroupOtherHouse.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserViktor.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupRavenclaw.Slug,
				GroupName: UserGroupRavenclaw.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserCho.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupRavenclaw.Slug,
				GroupName: UserGroupRavenclaw.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserDraco.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupSlytherin.Slug,
				GroupName: UserGroupSlytherin.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserSnape.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupSlytherin.Slug,
				GroupName: UserGroupSlytherin.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			{
				UserSlug: sql.NullString{
					String: UserLucius.Slug,
					Valid:  true,
				},
				GroupSlug: UserGroupSlytherin.Slug,
				GroupName: UserGroupSlytherin.Name,
				Deleted: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
		}

		result, err := services.SortUsersInToGroups(slugMap)
		require.NoError(t, err)
		require.Equal(t, int(5), len(result))

		require.Equal(t, UserGroupGryffindor.Name, result[0].Name)
		require.Equal(t, UserGroupGryffindor.Slug, result[0].Slug)
		require.Equal(t, false, result[0].Deleted)
		for _, userSlug := range result[0].UserSlugs {
			require.Contains(t, []string{UserHarry.Slug, UserGinny.Slug, UserRon.Slug, UserHermione.Slug}, userSlug)
		}

		require.Equal(t, UserGroupHufflepuff.Name, result[1].Name)
		require.Equal(t, UserGroupHufflepuff.Slug, result[1].Slug)
		require.Equal(t, false, result[1].Deleted)
		for _, userSlug := range result[1].UserSlugs {
			require.Contains(t, []string{UserFleur.Slug, UserCedric.Slug}, userSlug)
		}

		require.Equal(t, UserGroupOtherHouse.Name, result[2].Name)
		require.Equal(t, UserGroupOtherHouse.Slug, result[2].Slug)
		require.Equal(t, false, result[2].Deleted)
		for _, userSlug := range result[2].UserSlugs {
			require.Contains(t, []string{UserViktor.Slug, UserCho.Slug}, userSlug)
		}

		require.Equal(t, UserGroupRavenclaw.Name, result[3].Name)
		require.Equal(t, UserGroupRavenclaw.Slug, result[3].Slug)
		require.Equal(t, false, result[3].Deleted)
		for _, userSlug := range result[3].UserSlugs {
			require.Contains(t, []string{UserViktor.Slug, UserCho.Slug}, userSlug)
		}

		require.Equal(t, UserGroupSlytherin.Name, result[4].Name)
		require.Equal(t, UserGroupSlytherin.Name, result[4].Slug)
		require.Equal(t, false, result[4].Deleted)
		for _, userSlug := range result[4].UserSlugs {
			require.Contains(t, []string{UserDraco.Slug, UserSnape.Slug, UserLucius.Slug}, userSlug)
		}

		// if len(slugMap) == 0
		result, err = services.SortUsersInToGroups(services.SlugMap{})
		require.Equal(t, int(0), len(result))

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
