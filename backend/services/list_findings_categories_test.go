// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"fmt"
	"testing"

	"github.com/theparanoids/ashirt-server/backend/dtos"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

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
