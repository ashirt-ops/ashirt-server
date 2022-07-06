// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateFindingCategory(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := contextForUser(UserRon, db)

	// verify non-admin cannot add a new category
	newCategory := "Bogus Category"
	_, err := services.CreateFindingCategory(ctx, db, newCategory)
	require.Error(t, err)

	// verify that admins can create new categories
	ctx = contextForUser(UserDumbledore, db)
	newCategory = "Legitimate Category"
	createdFindingCategory, err := services.CreateFindingCategory(ctx, db, newCategory)
	require.NoError(t, err)
	require.Equal(t, newCategory, createdFindingCategory.Category)

	// verify that categories cannot be duplicated
	_, err = services.CreateFindingCategory(ctx, db, newCategory)
	require.Error(t, err)
}

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

func TestListFindingsCategories(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUserCtx := contextForUser(UserRon, db)

	// set up some helpers
	seedCategories := HarryPotterSeedData.FindingCategories
	activeSeedCategories := make([]models.FindingCategory, 0)

	for _, cat := range seedCategories {
		if cat.DeletedAt == nil {
			activeSeedCategories = append(activeSeedCategories, cat)
		}
	}

	activeCategories, err := services.ListFindingCategories(normalUserCtx, db, false)
	require.NoError(t, err)

	// verify all active categories exist (size is the same, and each active entry in seed)
	coercedCateories, ok := activeCategories.([]*dtos.FindingCategory)
	require.True(t, ok)
	requireFindingCategoriesAlign(t, activeSeedCategories, coercedCateories)

	allCategories, err := services.ListFindingCategories(normalUserCtx, db, true)
	require.NoError(t, err)

	// verify all categories exist (size is the same, and each entry exists in seed)
	coercedCateories, ok = allCategories.([]*dtos.FindingCategory)
	require.True(t, ok)
	requireFindingCategoriesAlign(t, seedCategories, coercedCateories)
}

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

func requireFindingCategoriesAlign(t *testing.T, modelList []models.FindingCategory, dtoList []*dtos.FindingCategory) {
	require.Equal(t, len(modelList), len(dtoList))
	for _, cat := range dtoList {
		matchItem := *cat
		found := false

		for _, item := range modelList {
			if item.ID == matchItem.ID {
				require.Equal(t, item.Category, matchItem.Category)
				require.Equal(t, matchItem.Deleted, (item.DeletedAt != nil))
				found = true
				break
			}
		}

		if found == false {
			require.Fail(t,
				fmt.Sprintf("No match found for ID: %v (category: %v)",
					matchItem.ID, matchItem.Category))
		}
	}
}
