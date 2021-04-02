// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestListFindingsCategories(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := contextForUser(UserRon, db)

	// verify that default entries are present
	categories, err := services.ListFindingCategories(ctx, db)
	require.NoError(t, err)
	require.NotEmpty(t, categories)
}
