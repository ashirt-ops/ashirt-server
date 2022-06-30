// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"

	sq "github.com/Masterminds/squirrel"
)

func TestDeleteEvidenceNoPropogate(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})
	memStore := createPopulatedMemStore(HarryPotterSeedData)

	masterEvidence := EviFlyingCar
	i := services.DeleteEvidenceInput{
		OperationSlug:            OpChamberOfSecrets.Slug,
		EvidenceUUID:             masterEvidence.UUID,
		DeleteAssociatedFindings: false,
	}
	// populate content store
	contentStoreKey := masterEvidence.UUID // seed data shares full and thumb key ids

	getAssociatedTagCount := makeDBRowCounter(t, db, "tag_evidence_map", "evidence_id=?", masterEvidence.ID)
	require.True(t, getAssociatedTagCount() > 0, "Database should have associated tags to delete")

	getEvidenceCount := makeDBRowCounter(t, db, "evidence", "uuid=?", i.EvidenceUUID)
	require.Equal(t, int64(1), getEvidenceCount(), "Database should have evidence to delete")

	err := services.DeleteEvidence(ctx, db, memStore, i)
	require.NoError(t, err)
	require.Equal(t, int64(0), getEvidenceCount(), "Database should have deleted the evidence")
	require.Equal(t, int64(0), getAssociatedTagCount(), "Database should have deleted associated tags")
	_, err = memStore.Read(contentStoreKey)
	require.Error(t, err)
}

func TestDeleteEvidenceWithPropogation(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})
	memStore := createPopulatedMemStore(HarryPotterSeedData)

	masterEvidence := EviDobby
	i := services.DeleteEvidenceInput{
		OperationSlug:            OpChamberOfSecrets.Slug,
		EvidenceUUID:             masterEvidence.UUID,
		DeleteAssociatedFindings: true,
	}
	getAssociatedTagCount := makeDBRowCounter(t, db, "tag_evidence_map", "evidence_id=?", masterEvidence.ID)
	require.True(t, getAssociatedTagCount() > 0, "Database should have associated tags to delete")

	getEvidenceCount := makeDBRowCounter(t, db, "evidence", "uuid=?", masterEvidence.UUID)
	require.Equal(t, int64(1), getEvidenceCount(), "Database should have evidence to delete")

	getMappedFindingCount := makeDBRowCounter(t, db, "evidence_finding_map", "evidence_id=?", masterEvidence.ID)
	require.True(t, getMappedFindingCount() > 0, "Database should have some mapped finding to delete")

	associatedFindingIDs := getAssociatedFindings(t, db, masterEvidence.ID)
	require.True(t, len(associatedFindingIDs) > 0, "Database should have some associated finding to delete")

	err := services.DeleteEvidence(ctx, db, memStore, i)
	require.NoError(t, err)
	require.Equal(t, int64(0), getEvidenceCount(), "Database should have deleted the evidence")
	require.Equal(t, int64(0), getAssociatedTagCount(), "Database should have deleted evidence-to-tags mappings")
	require.Equal(t, int64(0), getMappedFindingCount(), "Database should have deleted evidence-to-findings mappings")
	postDeleteFindingIDs := []int64{}
	db.Select(&postDeleteFindingIDs, sq.Select("id").From("findings").Where(sq.Eq{"id": associatedFindingIDs}))
	require.Equal(t, []int64{}, postDeleteFindingIDs, "Associated findings should be removed")
}

func getAssociatedFindings(t *testing.T, db *database.Connection, evidenceID int64) []int64 {
	query := sq.Select("finding_id").From("evidence_finding_map").
		Where(sq.Eq{"evidence_id": evidenceID})

	var rtn []int64
	err := db.Select(&rtn, query)

	require.Nil(t, err)
	return rtn
}
