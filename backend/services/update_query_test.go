// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestUpdateQuery(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	masterQuery := QuerySalazarsHier
	input := services.UpdateQueryInput{
		OperationSlug: masterOp.Slug,
		ID:            masterQuery.ID,
		Name:          "New Name",
		Query:         "New Query",
	}

	err := services.UpdateQuery(ctx, db, input)
	require.NoError(t, err)

	updatedQuery := getQueryByID(t, db, masterQuery.ID)

	require.NoError(t, err)
	require.Equal(t, input.Name, updatedQuery.Name)
	require.Equal(t, input.Query, updatedQuery.Query)
}
