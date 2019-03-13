// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

type userValidator func(*testing.T, UserOpPermJoinUser, *dtos.UserOperationRole)

func TestListUsersForOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	allUserOpRoles := getUsersWithRoleForOperationByOperationID(t, db, masterOp.ID)
	require.NotEqual(t, len(allUserOpRoles), 0, "Some users should be attached to this operation")

	input := services.ListUsersForOperationInput{
		OperationSlug: masterOp.Slug,
	}

	content, err := services.ListUsersForOperation(ctx, db, input)
	require.NoError(t, err)

	require.Equal(t, len(content), len(allUserOpRoles))
	validateUserSets(t, content, allUserOpRoles, validateUser)
}

func validateUser(t *testing.T, expected UserOpPermJoinUser, actual *dtos.UserOperationRole) {
	require.Equal(t, expected.Slug, actual.User.Slug)
	require.Equal(t, expected.FirstName, actual.User.FirstName)
	require.Equal(t, expected.LastName, actual.User.LastName)
	require.Equal(t, expected.Role, actual.Role)
}

func validateUserSets(t *testing.T, dtoSet []*dtos.UserOperationRole, dbSet []UserOpPermJoinUser, validate userValidator) {
	var expected *UserOpPermJoinUser = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.Slug == dtoItem.User.Slug {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validate(t, *expected, dtoItem)
	}
}
