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

	startingOp := OpChamberOfSecrets
	endingOp := OpSorcerersStone
	sourceEvidence := EviPetrifiedHermione //shares tags between the two operations

	input := services.MoveEvidenceInput{
		SourceOperationSlug: startingOp.Slug,
		TargetOperationSlug: endingOp.Slug,
		EvidenceUUID:        sourceEvidence.UUID,
	}

	// scenario 1: User present in both, cannot write dst [should fail]
	ctx := contextForUser(UserHermione, db)
	err := services.MoveEvidence(ctx, db, input)
	require.Error(t, err)

	// scenario 2: User present in both, cannot write src [should fail]
	ctx = contextForUser(UserSeamus, db)
	err = services.MoveEvidence(ctx, db, input)
	require.Error(t, err)

	// scenario 3: User present in src, cannot write dst [should fail]
	ctx = contextForUser(UserGinny, db)
	err = services.MoveEvidence(ctx, db, input)
	require.Error(t, err)

	// scenario 4: User present in dst, cannot write src [should fail]
	ctx = contextForUser(UserNeville, db)
	err = services.MoveEvidence(ctx, db, input)
	require.Error(t, err)

	// // scenario 5: User present in both, cannot write to both [should succeed]
	ctx = contextForUser(UserHarry, db)
	err = services.MoveEvidence(ctx, db, input)
	require.NoError(t, err)

	updatedEvidence := getEvidenceByUUID(t, db, sourceEvidence.UUID)
	require.Equal(t, updatedEvidence.OperationID, endingOp.ID)
	associatedTags := getTagIDsFromEvidenceID(t, db, updatedEvidence.ID)
	require.Equal(t, sorted(associatedTags), sorted([]int64{CommonTagWhoSS.ID, CommonTagWhatSS.ID}))
}
