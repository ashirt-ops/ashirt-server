// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestListOperations(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserDumbledore.ID, &policy.FullAccess{}) // by convention, this user should have admin access to all

	fullOps := getOperations(t, db)
	require.NotEqual(t, len(fullOps), 0, "Some number of operations should exist")

	ops, err := services.ListOperations(ctx, db)
	require.NoError(t, err)
	require.Equal(t, len(ops), len(fullOps))
	for _, op := range ops {
		var expected *models.Operation = nil
		for _, fOp := range fullOps {
			if fOp.ID == op.ID {
				expected = &fOp
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validateOp(t, *expected, op)
	}

	ctx = fullContext(UserDraco.ID, &policy.Deny{}) // user should have access to nothing
	ops, err = services.ListOperations(ctx, db)
	require.NoError(t, err)
	require.Equal(t, 0, len(ops))
}

func validateOp(t *testing.T, expected models.Operation, actual *dtos.Operation) {
	require.Equal(t, expected.Slug, actual.Slug, "Slugs should match")
	require.Equal(t, expected.Name, actual.Name, "Names should match")
	require.Equal(t, expected.Status, actual.Status, "Status should match")
}
