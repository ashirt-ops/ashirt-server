// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestMoveEvidence(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := contextForUser(UserRon, db)

	startingOp := OpChamberOfSecrets
	endingOp := OpSorcerersStone
	sourceEvidence := EviPetrifiedHermione //shares tags between the two operations

	input := services.MoveEvidenceInput{
		SourceOperationSlug: startingOp.Slug,
		TargetOperationSlug: endingOp.Slug,
		EvidenceUUID:        sourceEvidence.UUID,
	}

	// verify that Ron (cannot read endingOp) cannot determine tag differences
	err := services.MoveEvidence(ctx, db, input)
	require.Error(t, err)

	// verify that Harry (can write to both) can determine tag differences
	ctx = contextForUser(UserHarry, db)
	err = services.MoveEvidence(ctx, db, input)
	require.NoError(t, err)

	updatedEvidence := getEvidenceByUUID(t, db, sourceEvidence.UUID)
	require.Equal(t, updatedEvidence.OperationID, endingOp.ID)
	associatedTags := getTagIDsFromEvidenceID(t, db, updatedEvidence.ID)
	require.Equal(t, sorted(associatedTags), sorted([]int64{CommonTagWhoSS.ID, CommonTagWhatSS.ID}))
}
