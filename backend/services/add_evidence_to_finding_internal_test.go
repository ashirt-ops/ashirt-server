// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/database"
	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

func TestAddEvidenceToFindingInternalFunctions(t *testing.T) {
	db := internalTestDBSetup(t)
	goodOp, badOp := setupBasicTestOperation(t, db)

	testBatchAddEvidence(t, db, goodOp, badOp)
	testBatchRemoveEvidence(t, db, goodOp)
}

func getEvidenceIDs(t *testing.T, db *database.Connection, findingID int64) []int64 {
	var list []int64
	err := db.Select(&list, sq.Select("evidence_id").
		From("evidence_finding_map").
		Where(sq.Eq{"finding_id": findingID}).
		OrderBy("evidence_id ASC"))
	require.NoError(t, err)
	return list
}

func testBatchAddEvidence(t *testing.T, db *database.Connection, goodOp, badOp mockOperation) {
	findingID := goodOp.Findings[1].ID
	initialEviIDs := getEvidenceIDs(t, db, findingID)
	err := batchAddEvidenceToFinding(db, []string{}, goodOp.Op.ID, findingID)
	require.NoError(t, err)
	idsAfterEmptyAdd := getEvidenceIDs(t, db, findingID)
	require.Equal(t, initialEviIDs, idsAfterEmptyAdd)

	err = batchAddEvidenceToFinding(db, []string{badOp.Evidence[0].UUID}, goodOp.Op.ID, findingID)
	require.NotNil(t, err)

	err = batchAddEvidenceToFinding(db, []string{goodOp.Evidence[0].UUID}, goodOp.Op.ID, findingID)
	require.NoError(t, err)
	idsAfterSingleAdd := getEvidenceIDs(t, db, findingID)
	require.Equal(t, 1, len(idsAfterSingleAdd))
	require.Equal(t, goodOp.Evidence[0].ID, idsAfterSingleAdd[0])
}

func testBatchRemoveEvidence(t *testing.T, db *database.Connection, goodOp mockOperation) {
	findingID := goodOp.Findings[0].ID
	_, err := db.Insert("evidence_finding_map", map[string]interface{}{"evidence_id": goodOp.Evidence[0].ID, "finding_id": findingID})
	require.NoError(t, err)
	_, err = db.Insert("evidence_finding_map", map[string]interface{}{"evidence_id": goodOp.Evidence[1].ID, "finding_id": findingID})
	require.NoError(t, err)

	initialEviIDs := getEvidenceIDs(t, db, findingID)
	err = batchRemoveEvidenceFromFinding(db, []string{}, findingID)
	require.NoError(t, err)
	idsAfterEmptyDelete := getEvidenceIDs(t, db, findingID)
	require.Equal(t, initialEviIDs, idsAfterEmptyDelete)

	err = batchRemoveEvidenceFromFinding(db, []string{goodOp.Evidence[0].UUID}, findingID)
	require.NoError(t, err)
	idsAfterSemiDelete := getEvidenceIDs(t, db, findingID)
	require.Equal(t, 1, len(idsAfterSemiDelete))
	require.Equal(t, goodOp.Evidence[1].ID, idsAfterSemiDelete[0])
}
