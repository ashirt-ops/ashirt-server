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
