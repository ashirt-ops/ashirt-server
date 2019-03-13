// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/models"

	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestUpdateOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	// tests for common fields
	masterOp := OpChamberOfSecrets
	input := services.UpdateOperationInput{
		OperationSlug: masterOp.Slug,
		Name:          "New Name",
		Status:        models.OperationStatusComplete,
	}
	require.NotEqual(t, masterOp.Status, input.Status)

	err := services.UpdateOperation(ctx, db, input)
	require.NoError(t, err)
	updatedOperation, err := services.ReadOperation(ctx, db, masterOp.Slug)
	require.NoError(t, err)
	require.Equal(t, input.Name, updatedOperation.Name)
	require.Equal(t, input.Status, updatedOperation.Status)
}
