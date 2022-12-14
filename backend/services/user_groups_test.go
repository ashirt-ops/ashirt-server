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

func GetUserIDsFromGroup(db *database.Connection, groupName string) ([]int64, error) {
	var userGroupId int64
	err := db.Get(&userGroupId, sq.Select("id").
		From("user_groups").
		Where(sq.Eq{
			"slug": groupName,
		}))
	if err != nil {
		s := fmt.Sprintf("Cannot get user group id for group %q", groupName)
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

func TestCreateAndDeleteUserGroup(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		name := "testGroup"
		i := services.CreateUserGroupInput{
			Name: name,
			UserSlugs: []string{
				UserRon.Slug,
				UserAlastor.Slug,
				UserHagrid.Slug,
			},
		}

		adminUser := UserDumbledore
		ctx := contextForUser(adminUser, db)
		_, err := services.CreateUserGroup(ctx, db, i)
		require.NoError(t, err)

		userIDs, err := GetUserIDsFromGroup(db, name)
		require.NoError(t, err)
		require.Equal(t, 3, len(userIDs))
		for _, userID := range userIDs {
			require.Contains(t, []int64{UserRon.ID, UserAlastor.ID, UserHagrid.ID}, userID)
		}

		_, err = services.CreateUserGroup(ctx, db, i)
		assert.ErrorContains(t, err, "Unable to create user group. User group slug already exists")

		err = services.DeleteUserGroup(ctx, db, name)
		require.NoError(t, err)

		userIDs, err = GetUserIDsFromGroup(db, name)
		require.NoError(t, err)
		require.Equal(t, 0, len(userIDs))
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
		require.Equal(t, result.PageNumber, int64(1))
		require.Equal(t, result.PageSize, int64(10))
		require.Equal(t, result.TotalCount, int64(4))
		require.Equal(t, len(usergroups), 5)
		require.NoError(t, err)
	})
}
