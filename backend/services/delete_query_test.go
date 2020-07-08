// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteQuery(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	i := services.DeleteQueryInput{
		OperationSlug: OpChamberOfSecrets.Slug,
		ID:            QuerySalazarsHier.ID,
	}

	getQueryCount := makeDBRowCounter(t, db, "queries", "id=?", i.ID)
	require.Equal(t, int64(1), getQueryCount(), "Database should have item to delete")

	err := services.DeleteQuery(ctx, db, i)
	require.NoError(t, err)
	require.Equal(t, int64(0), getQueryCount(), "Database should have deleted the item")
}
