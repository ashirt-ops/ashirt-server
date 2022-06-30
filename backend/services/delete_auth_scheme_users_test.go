// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	localConsts "github.com/theparanoids/ashirt-server/backend/authschemes/localauth/constants"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteAuthSchemeUsers(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)
	schemeName := localConsts.Code

	baseUsers := getUsersForAuth(t, db, schemeName)
	require.Greater(t, len(baseUsers), 0)

	// verify non-admins have no access
	err := services.DeleteAuthSchemeUsers(ctx, db, schemeName)
	require.Error(t, err)

	// verify admins have access + effect works
	ctx = simpleFullContext(adminUser)
	err = services.DeleteAuthSchemeUsers(ctx, db, schemeName)
	require.NoError(t, err)

	updatedUsers := getUsersForAuth(t, db, schemeName)
	require.Equal(t, 0, len(updatedUsers))

	// verify admins cannot delete recovery
	err = services.DeleteAuthSchemeUsers(ctx, db, recoveryConsts.Code)
	require.Error(t, err)
}
