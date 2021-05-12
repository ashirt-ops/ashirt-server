// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestUpdateFindingCategory(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := contextForUser(UserRon, db)

	targetCategory := ProductFindingCategory

	// verify non-admin cannot add a new category
	i := services.UpdateFindingCategoryInput{
		Category: "My new category",
		ID:       targetCategory.ID,
	}
	err := services.UpdateFindingCategory(ctx, db, i)
	require.Error(t, err)

	// verify that admins can create new categories
	ctx = contextForUser(UserDumbledore, db)
	err = services.UpdateFindingCategory(ctx, db, i)
	require.NoError(t, err)

	updatedFindingList := make([]models.FindingCategory, 0)
	for _, item := range HarryPotterSeedData.FindingCategories {
		if item == targetCategory {
			updatedCategory := targetCategory
			updatedCategory.Category = i.Category
			updatedFindingList = append(updatedFindingList, updatedCategory)
		} else {
			updatedFindingList = append(updatedFindingList, item)
		}
	}

	allCategories, err := services.ListFindingCategories(ctx, db, true)
	require.NoError(t, err)
	coercedCateories, ok := allCategories.([]*dtos.FindingCategory)
	require.True(t, ok)

	requireFindingCategoriesAlign(t, updatedFindingList, coercedCateories)
}
