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

func TestDeleteFindingCategory(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := contextForUser(UserRon, db)
	deleteTargetCategory := ProductFindingCategory
	restoreTargetCategory := DeletedCategory

	deleteInput := services.DeleteFindingCategoryInput{
		DoDelete:          true,
		FindingCategoryId: deleteTargetCategory.ID,
	}

	// verify that normal users cannot delete categories
	err := services.DeleteFindingCategory(ctx, db, deleteInput)
	require.Error(t, err)

	// verify that admins can delete categories
	ctx = contextForUser(UserDumbledore, db)
	err = services.DeleteFindingCategory(ctx, db, deleteInput)
	require.NoError(t, err)

	// verify that admins can restore categories
	restoreInput := services.DeleteFindingCategoryInput{
		DoDelete:          false,
		FindingCategoryId: restoreTargetCategory.ID,
	}

	// verify that categories cannot be duplicated
	err = services.DeleteFindingCategory(ctx, db, restoreInput)
	require.NoError(t, err)

	// check list results
	updatedFindingList := make([]models.FindingCategory, 0)
	for _, item := range HarryPotterSeedData.FindingCategories {
		if item == deleteTargetCategory {
			deletedTime := GetInternalClock().Now()
			deleteTargetCategory.DeletedAt = &deletedTime
			updatedFindingList = append(updatedFindingList, deleteTargetCategory)
		} else if item == restoreTargetCategory {
			restoreTargetCategory.DeletedAt = nil
			updatedFindingList = append(updatedFindingList, restoreTargetCategory)
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
