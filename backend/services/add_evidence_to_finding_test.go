// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestAddEvidenceToFinding(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	masterFinding := FindingBook2Magic
	evidenceToAdd1 := EviSpiderAragog
	evidenceToAdd2 := EviMoaningMyrtle
	evidenceToRemove1 := EviDobby
	evidenceToRemove2 := EviFlyingCar

	initialEvidenceList := getEvidenceIDsFromFinding(t, db, masterFinding.ID)

	expectedEvidenceSet := make(map[int64]bool)
	for _, id := range initialEvidenceList {
		if id != evidenceToRemove1.ID && id != evidenceToRemove2.ID {
			expectedEvidenceSet[id] = true
		}
	}
	expectedEvidenceSet[evidenceToAdd1.ID] = true
	expectedEvidenceSet[evidenceToAdd2.ID] = true
	expectedEvidenceList := make([]int64, 0, len(expectedEvidenceSet))
	for key, v := range expectedEvidenceSet {
		if v {
			expectedEvidenceList = append(expectedEvidenceList, key)
		}
	}

	i := services.AddEvidenceToFindingInput{
		OperationSlug:    masterOp.Slug,
		FindingUUID:      masterFinding.UUID,
		EvidenceToAdd:    []string{evidenceToAdd1.UUID, evidenceToAdd2.UUID},
		EvidenceToRemove: []string{evidenceToRemove1.UUID, evidenceToRemove2.UUID},
	}
	err := services.AddEvidenceToFinding(ctx, db, i)
	require.NoError(t, err)

	changedEvidenceList := getEvidenceIDsFromFinding(t, db, masterFinding.ID)

	require.Equal(t, sorted(expectedEvidenceList), sorted(changedEvidenceList))
}
