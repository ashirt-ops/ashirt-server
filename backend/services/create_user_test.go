// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateUser(t *testing.T) {
	db := initTest(t)

	// verify first user is an admin
	i := services.CreateUserInput{
		FirstName: "Luna",
		LastName:  "Lovegood",
		Slug:      "luna.lovegood",
		Email:     "luna.lovegood@hogwarts.edu",
	}

	createUserOutput, err := services.CreateUser(db, i)
	require.NoError(t, err)
	require.Equal(t, createUserOutput.RealSlug, i.Slug)
	luna := getUserProfile(t, db, createUserOutput.UserID)

	require.Equal(t, true, luna.Admin)
	require.Equal(t, luna.FirstName, i.FirstName)
	require.Equal(t, luna.Email, i.Email)
	require.Equal(t, luna.LastName, i.LastName)

	// Verify re-register will fail (due to unique email constraint)
	createUserOutput, err = services.CreateUser(db, i)
	require.Error(t, err)

	// Verify 2nd user (non-admin, no matching slug)
	i.Email = "luna.lovegood+extra@hogwarts.edu" // change the password to something that won't exist
	createUserOutput, err = services.CreateUser(db, i)
	require.NoError(t, err)
	// Since Luna's already exists, a new slug should be created
	require.NotEqual(t, i.Slug, createUserOutput.RealSlug)
	require.Contains(t, createUserOutput.RealSlug, i.Slug)
	newLuna := getUserProfile(t, db, createUserOutput.UserID)

	require.Equal(t, false, newLuna.Admin)
	require.Equal(t, i.FirstName, newLuna.FirstName)
	require.Equal(t, i.Email, newLuna.Email)
	require.Equal(t, i.LastName, newLuna.LastName)
}

func TestCreateHeadlessUser(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	i := services.CreateUserInput{
		FirstName: "Extra",
		LastName:  "Headless Hunt Member",
		Slug:      "sir.nobody",
		Email:     "sir.nobody@hogwarts.edu",
	}

	// Verify non-admin can not create headless users
	ctx := simpleFullContext(UserHarry)
	_, err := services.CreateHeadlessUser(ctx, db, i)
	require.Error(t, err)

	ctx = simpleFullContext(UserDumbledore)
	result, err := services.CreateHeadlessUser(ctx, db, i)
	require.NoError(t, err)

	foundUser := getUserBySlug(t, db, result.RealSlug)

	require.True(t, foundUser.Headless)
}
