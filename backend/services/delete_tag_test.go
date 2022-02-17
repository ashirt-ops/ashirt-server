// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteTag(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	op := OpChamberOfSecrets
	i := services.DeleteTagInput{
		ID:            TagEarth.ID,
		OperationSlug: op.Slug,
	}

	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	err := services.DeleteTag(ctx, db, i)
	require.NoError(t, err)

	require.NotContains(t, getTagFromOperationID(t, db, op.ID), TagEarth, "TagEarth should have been deleted")
}

func TestDeleteDefaultTag(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	tagToRemove := DefaultTagWho
	normalUser := UserRon
	adminUser := UserDumbledore

	i := services.DeleteDefaultTagInput{
		ID: tagToRemove.ID,
	}

	// verify that a normal user cannot delete default tags
	ctx := simpleFullContext(normalUser)
	err := services.DeleteDefaultTag(ctx, db, i)
	require.Error(t, err)

	// verify that an admin can delete default tags
	ctx = simpleFullContext(adminUser)
	err = services.DeleteDefaultTag(ctx, db, i)
	require.NoError(t, err)
	require.NotContains(t, getDefaultTags(t, db), tagToRemove)
}
