// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteTag(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	op := OpChamberOfSecrets
	i := services.DeleteTagInput{
		ID:            TagEarth.ID,
		OperationSlug: op.Slug,
	}

	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	err := services.DeleteTag(ctx, db, i)
	require.NoError(t, err)

	require.NotContains(t, getTagFromOperationID(t, db, op.ID), TagEarth, "TagEarth should have been deleted")
}
