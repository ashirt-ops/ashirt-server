// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateTag(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)

	op := OpSorcerersStone
	i := services.CreateTagInput{
		Name:          "New Tag",
		ColorName:     "indigo",
		OperationSlug: op.Slug,
	}

	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	createdTag, err := services.CreateTag(ctx, db, i)
	require.NoError(t, err)
	require.Equal(t, createdTag.Name, i.Name)
	require.NotContains(t, HarryPotterSeedData.AllInitialTagIds(), createdTag.ID, "Should have new ID")

	updatedTag := getTagByID(t, db, createdTag.ID)

	require.Equal(t, op.ID, updatedTag.OperationID, "is in right operation")
}

func TestCreateDefaultTag(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	adminUser := UserDumbledore

	i := services.CreateDefaultTagInput{
		Name:      "New Tag",
		ColorName: "indigo",
	}

	// verify that a normal cannot create a new default tag
	ctx := simpleFullContext(normalUser)
	_, err := services.CreateDefaultTag(ctx, db, i)
	require.Error(t, err)

	// verify that an admin can create a new default tag
	ctx = simpleFullContext(adminUser)
	createdTag, err := services.CreateDefaultTag(ctx, db, i)
	require.NoError(t, err)
	require.Equal(t, createdTag.Name, i.Name)
	require.NotContains(t, HarryPotterSeedData.AllInitialDefaultTagIds(), createdTag.ID, "Should have new ID")
}
