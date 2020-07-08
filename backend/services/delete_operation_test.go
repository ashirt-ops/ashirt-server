// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestDeleteOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserHarry.ID, &policy.Deny{})
	memStore := createPopulatedMemStore(HarryPotterSeedData)

	masterOp := OpChamberOfSecrets
	originalEvidence := getEvidenceForOperation(t, db, masterOp.ID)

	// Verify that non-admins cannot delete
	err := services.DeleteOperation(ctx, db, memStore, masterOp.Slug)
	require.Error(t, err)

	// Verify admins can delete
	ctx = fullContext(UserRon.ID, &policy.FullAccess{})
	err = services.DeleteOperation(ctx, db, memStore, masterOp.Slug)
	require.NoError(t, err)
	// ensure content was removed
	for _, evi := range originalEvidence {
		_, err = memStore.Read(evi.FullImageKey)
		require.Error(t, err)
		_, err = memStore.Read(evi.ThumbImageKey)
		require.Error(t, err)
	}
	var dbOp models.Operation
	err = db.Get(&dbOp, sq.Select("*").From("operations").Where(sq.Eq{"id": masterOp.ID}))
	// assuming that if this row was deleted, then all other rows must have been deleted (via foreign key constraint)
	require.Error(t, err)

	// Verify Super admins can delete
	// TODO
}
