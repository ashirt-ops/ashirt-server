// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
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
