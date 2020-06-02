// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestListQueriesForOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	allQueries := getQueriesForOperationID(t, db, masterOp.ID)
	require.NotEqual(t, len(allQueries), 0, "Some number of queries should exist")

	foundQueries, err := services.ListQueriesForOperation(ctx, db, masterOp.Slug)
	require.NoError(t, err)
	require.Equal(t, len(foundQueries), len(allQueries))
	validateQuerySets(t, foundQueries, allQueries, validateQuery)
}

type queryValidator func(*testing.T, models.Query, *dtos.Query)

func validateQuery(t *testing.T, expected models.Query, actual *dtos.Query) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Query, actual.Query)
	require.Equal(t, expected.Type, actual.Type)
}

func validateQuerySets(t *testing.T, dtoSet []*dtos.Query, dbSet []models.Query, validator queryValidator) {
	var expected *models.Query = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.ID == dtoItem.ID {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validator(t, *expected, dtoItem)
	}
}
