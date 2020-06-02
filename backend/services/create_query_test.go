// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestCreateQuery(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	op := OpChamberOfSecrets
	i := services.CreateQueryInput{
		OperationSlug: op.Slug,
		Name:          "Evidence By author",
		Query:         "<query goes here>",
		Type:          "findings",
	}
	createdQuery, err := services.CreateQuery(ctx, db, i)
	require.NoError(t, err)
	fullQuery := getQueryByID(t, db, createdQuery.ID)
	require.Equal(t, i.Name, fullQuery.Name)
	require.Equal(t, i.Type, fullQuery.Type)
	require.Equal(t, i.Query, fullQuery.Query)
}
