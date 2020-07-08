// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	targetUser := UserRon
	admin := UserDumbledore

	require.True(t, 0 < countRows(t, db, "api_keys", "user_id=?", targetUser.ID))
	require.True(t, 0 < countRows(t, db, "auth_scheme_data", "user_id=?", targetUser.ID))
	require.True(t, 0 < countRows(t, db, "user_operation_permissions", "user_id=?", targetUser.ID))

	// verify that non-admins cannot delete
	ctx := fullContext(UserDraco.ID, &policy.FullAccess{})
	err := services.DeleteUser(ctx, db, targetUser.Slug)
	require.Error(t, err)

	// verify user cannot delete themselves
	ctx = fullContextAsAdmin(admin.ID, &policy.FullAccess{})
	err = services.DeleteUser(ctx, db, admin.Slug)
	require.NotNil(t, err)

	// Verify delete actually works
	err = services.DeleteUser(ctx, db, targetUser.Slug)
	require.Nil(t, err)

	require.True(t, 0 == countRows(t, db, "api_keys", "user_id=?", targetUser.ID))
	require.True(t, 0 == countRows(t, db, "auth_scheme_data", "user_id=?", targetUser.ID))
	require.True(t, 0 == countRows(t, db, "user_operation_permissions", "user_id=?", targetUser.ID))

	var user models.User
	err = db.Get(&user, sq.Select("*").From("users").Where(sq.Eq{"id": targetUser.ID}))
	require.Nil(t, err)
	require.NotNil(t, user.DeletedAt)
}
