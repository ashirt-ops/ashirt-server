// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"

	sq "github.com/Masterminds/squirrel"
)

func TestSetUserOperationRole(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	targetUser := UserHarry
	targetRole := policy.OperationRoleRead
	input := services.SetUserOperationRoleInput{
		OperationSlug: masterOp.Slug,
		UserSlug:      targetUser.Slug,
		Role:          targetRole,
	}

	initialRole := HarryPotterSeedData.UserRoleForOp(targetUser, masterOp)
	require.NotContains(t, []policy.OperationRole{targetRole, ""}, initialRole, "Test user should both have a role, but not have the role we want to use")

	err := services.SetUserOperationRole(ctx, db, input)
	require.NoError(t, err)

	getDBRole := func() (string, error) {
		var newRole string
		err := db.Get(&newRole, sq.Select("role").
			From("user_operation_permissions").
			Where(sq.Eq{"operation_id": masterOp.ID, "user_id": targetUser.ID}))
		return newRole, err
	}
	newRole, err := getDBRole()
	require.NoError(t, err)
	require.Equal(t, string(targetRole), newRole)

	input = services.SetUserOperationRoleInput{
		OperationSlug: masterOp.Slug,
		UserSlug:      targetUser.Slug,
		Role:          "",
	}

	err = services.SetUserOperationRole(ctx, db, input)
	require.NoError(t, err)

	_, err = getDBRole()
	require.True(t, database.IsEmptyResultSetError(err))

	targetRole = policy.OperationRoleAdmin
	input = services.SetUserOperationRoleInput{
		OperationSlug: masterOp.Slug,
		UserSlug:      targetUser.Slug,
		Role:          targetRole,
	}
	err = services.SetUserOperationRole(ctx, db, input)
	require.NoError(t, err)

	newRole, err = getDBRole()
	require.NoError(t, err)
	require.Equal(t, string(targetRole), newRole)
}
