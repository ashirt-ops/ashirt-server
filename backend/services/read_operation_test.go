// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestReadOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets

	retrievedOp, err := services.ReadOperation(ctx, db, masterOp.Slug)
	require.NoError(t, err)

	require.Equal(t, masterOp.Slug, retrievedOp.Slug)
	require.Equal(t, masterOp.Name, retrievedOp.Name)
	require.Equal(t, masterOp.Status, retrievedOp.Status)
	require.Equal(t, len(HarryPotterSeedData.UsersForOp(masterOp)), retrievedOp.NumUsers)
}
