// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestListOperationsForAdmin(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContextAsAdmin(UserDumbledore.ID, &policy.FullAccess{})

	fullOps := getOperations(t, db)
	require.NotEqual(t, len(fullOps), 0, "Some number of operations should exist")

	ops, err := services.ListOperationsForAdmin(ctx, db)
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

	// verify non admins don't have access

	ctx = fullContext(UserDraco.ID, &policy.FullAccess{}) // Note: not an admin
	ops, err = services.ListOperationsForAdmin(ctx, db)
	require.Error(t, err)
	require.Equal(t, "Requesting user is not an admin", err.Error())
}
