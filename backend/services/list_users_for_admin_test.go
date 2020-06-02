// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestListUsersForAdmin(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	allUsers := getAllUsers(t, db)
	allDeletedUsers := getAllDeletedUsers(t, db)

	input := services.ListUsersForAdminInput{
		Pagination: services.Pagination{
			Page:     1,
			PageSize: 250,
		},
		IncludeDeleted: false,
	}
	input.Pagination.SetMaxItems(input.Pagination.PageSize) // force constrain not to affect us

	// verify access restricted for non-admins
	ctx := fullContext(UserDraco.ID, &policy.FullAccess{}) // Note: not an admin
	_, err := services.ListUsersForAdmin(ctx, db, input)
	require.Error(t, err)
	require.Equal(t, "Requesting user is not an admin", err.Error())

	// Verify admins can list users (no deleted users)
	ctx = fullContextAsAdmin(UserDumbledore.ID, &policy.FullAccess{})
	pagedUsers, err := services.ListUsersForAdmin(ctx, db, input)
	require.NoError(t, err)

	require.Equal(t, input.Pagination.PageSize, pagedUsers.PageSize)
	require.Equal(t, int64(len(allUsers)-len(allDeletedUsers)), pagedUsers.TotalCount)

	usersDto, ok := pagedUsers.Content.([]*dtos.UserAdminView)
	require.True(t, ok)
	dtoIndex := 0
	for i := 0; i < len(allUsers); i++ {
		if allUsers[i].DeletedAt != nil {
			continue
		}
		require.Equal(t, allUsers[i].Slug, usersDto[dtoIndex].Slug)
		require.Equal(t, allUsers[i].FirstName, usersDto[dtoIndex].FirstName)
		require.Equal(t, allUsers[i].LastName, usersDto[dtoIndex].LastName)
		require.Equal(t, allUsers[i].Admin, usersDto[dtoIndex].Admin)
		require.Equal(t, allUsers[i].Disabled, usersDto[dtoIndex].Disabled)
		require.Equal(t, allUsers[i].Email, usersDto[dtoIndex].Email)
		require.Equal(t, allUsers[i].Headless, usersDto[dtoIndex].Headless)
		require.Equal(t, false, usersDto[dtoIndex].Deleted)
		dtoIndex++
	}

	// verify deleted users can be shown
	input.IncludeDeleted = true
	pagedUsers, err = services.ListUsersForAdmin(ctx, db, input)
	require.Nil(t, err)

	usersDto, _ = pagedUsers.Content.([]*dtos.UserAdminView)
	for i := 0; i < len(allUsers); i++ {
		require.Equal(t, allUsers[i].Slug, usersDto[i].Slug)
		require.Equal(t, allUsers[i].FirstName, usersDto[i].FirstName)
		require.Equal(t, allUsers[i].LastName, usersDto[i].LastName)
		require.Equal(t, allUsers[i].Admin, usersDto[i].Admin)
		require.Equal(t, allUsers[i].Disabled, usersDto[i].Disabled)
		require.Equal(t, allUsers[i].Email, usersDto[i].Email)
		require.Equal(t, allUsers[i].Headless, usersDto[i].Headless)
		require.Equal(t, (allUsers[i].DeletedAt != nil), usersDto[i].Deleted)
	}
}
