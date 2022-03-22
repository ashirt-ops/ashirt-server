// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestListOperations(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	validateOperationList := func(receivedOps []*dtos.Operation, expectedOps []models.Operation) {
		for _, op := range receivedOps {
			var expected *models.Operation = nil
			for _, fOp := range expectedOps {
				if fOp.Slug == op.Slug {
					expected = &fOp
					break
				}
			}
			require.NotNil(t, expected, "Result should have matching value")
			validateOp(t, *expected, op)
		}
	}

	normalUser := UserRon
	expectedOps := getOperationsForUser(t, db, normalUser.ID)

	ops, err := services.ListOperations(contextForUser(normalUser, db), db)
	require.NoError(t, err)
	require.Equal(t, len(ops), len(expectedOps))
	validateOperationList(ops, expectedOps)

	// validate headless users
	headlessUser := UserHeadlessNick
	fullOps := getOperations(t, db)

	ops, err = services.ListOperations(contextForUser(headlessUser, db), db)
	require.NoError(t, err)
	require.Equal(t, len(ops), len(fullOps))
	validateOperationList(ops, fullOps)
}

func validateOp(t *testing.T, expected models.Operation, actual *dtos.Operation) {
	require.Equal(t, expected.Slug, actual.Slug, "Slugs should match")
	require.Equal(t, expected.Name, actual.Name, "Names should match")
	require.Equal(t, expected.Status, actual.Status, "Status should match")
}
