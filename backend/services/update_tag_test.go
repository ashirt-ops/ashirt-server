// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestUpdateTag(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	op := OpChamberOfSecrets
	i := services.UpdateTagInput{
		ID:            TagEarth.ID,
		OperationSlug: op.Slug,
		Name:          "Moon",
		ColorName:     "green",
	}

	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	err := services.UpdateTag(ctx, db, i)
	require.NoError(t, err)

	updatedTag := getTagByID(t, db, TagEarth.ID)
	require.Equal(t, models.Tag{
		ID:          TagEarth.ID,
		OperationID: op.ID,
		Name:        "Moon",
		ColorName:   "green",
		CreatedAt:   TagEarth.CreatedAt,
		UpdatedAt:   updatedTag.UpdatedAt,
	}, updatedTag)
}

func TestUpdateDefaultTag(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	adminUser := UserDumbledore
	tagToUpdate := DefaultTagWho

	i := services.UpdateDefaultTagInput{
		ID:        tagToUpdate.ID,
		Name:      "How",
		ColorName: "green",
	}

	// verify that a normal user cannot update a default tags
	ctx := simpleFullContext(normalUser)
	err := services.UpdateDefaultTag(ctx, db, i)
	require.Error(t, err)

	// verify that an admin can update default tags
	ctx = simpleFullContext(adminUser)
	err = services.UpdateDefaultTag(ctx, db, i)
	require.NoError(t, err)

	updatedTag := getDefaultTagByID(t, db, tagToUpdate.ID)
	require.Equal(t, models.DefaultTag{
		ID:        tagToUpdate.ID,
		Name:      i.Name,
		ColorName: i.ColorName,
		CreatedAt: tagToUpdate.CreatedAt,
		UpdatedAt: updatedTag.UpdatedAt,
	}, updatedTag)
}
