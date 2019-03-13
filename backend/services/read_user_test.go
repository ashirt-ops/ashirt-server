// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestReadUser(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	targetUser := UserHarry
	adminUser := UserDumbledore
	ctx := simpleFullContext(normalUser)

	// verify read-self
	retrievedUser, err := services.ReadUser(ctx, db, "")
	require.Nil(t, err)
	verifyRetrievedUser(t, normalUser, retrievedUser)

	// verify read-self alternative (userslug provided)
	retrievedUser, err = services.ReadUser(ctx, db, normalUser.Slug)
	require.Nil(t, err)
	verifyRetrievedUser(t, normalUser, retrievedUser)

	// verify read-other (non-admin : should fail)
	_, err = services.ReadUser(ctx, db, targetUser.Slug)
	require.NotNil(t, err)

	// verify read-other (as admin)
	ctx = simpleFullContext(adminUser)
	retrievedUser, err = services.ReadUser(ctx, db, targetUser.Slug)
	require.Nil(t, err)
	verifyRetrievedUser(t, targetUser, retrievedUser)
}

func verifyRetrievedUser(t *testing.T, expectedUser models.User, retrievedUser *dtos.UserOwnView) {
	require.Equal(t, expectedUser.Slug, retrievedUser.Slug)
	require.Equal(t, expectedUser.FirstName, retrievedUser.FirstName)
	require.Equal(t, expectedUser.LastName, retrievedUser.LastName)
	require.Equal(t, expectedUser.Email, retrievedUser.Email)
}
