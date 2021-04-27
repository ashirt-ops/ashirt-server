// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateFindingCategory(t *testing.T) {
	db := initTest(t)
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
