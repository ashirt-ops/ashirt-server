// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

func TestEnsureTagIDsBelongToOperation(t *testing.T) {
	db := internalTestDBSetup(t)
	goodOp, badOp := setupBasicTestOperation(t, db)

	opCopy := opDtoToModel(goodOp.Op, goodOp.ID)

	// Add a tag for testing
	newTags := []models.Tag{
		{OperationID: opCopy.ID, Name: "good tag", ColorName: "black"},
		{OperationID: badOp.ID, Name: "bad tag", ColorName: "black"},
	}
	err := db.BatchInsert("tags", len(newTags), func(i int) map[string]interface{} {
		return map[string]interface{}{
			"operation_id": newTags[i].OperationID,
			"color_name":   newTags[i].ColorName,
			"name":         newTags[i].Name,
		}
	})
	require.NoError(t, err)
	var goodTagID int64
	getIDQuery := sq.Select("id").From("tags").Limit(1)

	err = db.Get(&goodTagID, getIDQuery.Where(sq.Eq{"operation_id": opCopy.ID}))
	require.NoError(t, err)

	var badTagID int64
	err = db.Get(&badTagID, getIDQuery.Where(sq.Eq{"operation_id": badOp.ID}))
	require.NoError(t, err)

	// No-data check (should be a noop)
	err = ensureTagIDsBelongToOperation(db, []int64{}, &opCopy)
	require.NoError(t, err)

	// Only non-existant tag check
	err = ensureTagIDsBelongToOperation(db, []int64{badTagID}, &opCopy)
	require.NotNil(t, err)

	// Has-tag check
	err = ensureTagIDsBelongToOperation(db, []int64{goodTagID}, &opCopy)
	require.NoError(t, err)

	// Has-both check
	err = ensureTagIDsBelongToOperation(db, []int64{goodTagID, badTagID}, &opCopy)
	require.NotNil(t, err)
}

func TestLookupOperation(t *testing.T) {
	db := internalTestDBSetup(t)
	goodOp, _ := setupBasicTestOperation(t, db)

	lookedUp, err := lookupOperation(db, goodOp.Op.Slug)
	require.NoError(t, err)
	require.Equal(t, goodOp.Op.Name, lookedUp.Name)
	require.Equal(t, goodOp.ID, lookedUp.ID)
	require.Equal(t, goodOp.Op.Status, lookedUp.Status)

	lookedUp, err = lookupOperation(db, "not-a-slug")
	require.Error(t, err)
	require.Equal(t, &models.Operation{}, lookedUp)
}

func TestLookupOperationFinding(t *testing.T) {
	db := internalTestDBSetup(t)
	goodOp, badOp := setupBasicTestOperation(t, db)

	foundOp, foundFinding, err := lookupOperationFinding(db, goodOp.Op.Slug, goodOp.Findings[0].UUID)
	require.NoError(t, err)
	require.Equal(t, goodOp.Op.Name, foundOp.Name)
	require.Equal(t, goodOp.ID, foundOp.ID)
	require.Equal(t, goodOp.Op.Status, foundOp.Status)

	require.Equal(t, goodOp.Findings[0].ID, foundFinding.ID)
	require.Equal(t, goodOp.Findings[0].UUID, foundFinding.UUID)
	require.Equal(t, goodOp.Findings[0].OperationID, foundFinding.OperationID)
	require.Equal(t, goodOp.Findings[0].CategoryID, foundFinding.CategoryID)
	require.Equal(t, goodOp.Findings[0].Title, foundFinding.Title)
	require.Equal(t, goodOp.Findings[0].Description, foundFinding.Description)

	_, _, err = lookupOperationFinding(db, "not-a-slug", goodOp.Findings[0].UUID)
	require.Error(t, err)

	_, _, err = lookupOperationFinding(db, goodOp.Op.Slug, "not-a-uuid")
	require.Error(t, err)

	_, _, err = lookupOperationFinding(db, badOp.Op.Slug, goodOp.Findings[0].UUID)
	require.Error(t, err)
}

func TestLookupOperationEvidence(t *testing.T) {
	db := internalTestDBSetup(t)
	goodOp, badOp := setupBasicTestOperation(t, db)

	foundOp, foundEvidence, err := lookupOperationEvidence(db, goodOp.Op.Slug, goodOp.Evidence[0].UUID)
	require.NoError(t, err)
	require.Equal(t, goodOp.Op.Name, foundOp.Name)
	require.Equal(t, goodOp.ID, foundOp.ID)
	require.Equal(t, goodOp.Op.Status, foundOp.Status)

	require.Equal(t, goodOp.Evidence[0].ID, foundEvidence.ID)
	require.Equal(t, goodOp.Evidence[0].UUID, foundEvidence.UUID)
	require.Equal(t, goodOp.Evidence[0].OperationID, foundEvidence.OperationID)
	require.Equal(t, goodOp.Evidence[0].OperatorID, foundEvidence.OperatorID)
	require.Equal(t, goodOp.Evidence[0].Description, foundEvidence.Description)
	require.Equal(t, goodOp.Evidence[0].ContentType, foundEvidence.ContentType)

	_, _, err = lookupOperationEvidence(db, "not-a-slug", goodOp.Evidence[0].UUID)
	require.Error(t, err)

	_, _, err = lookupOperationEvidence(db, goodOp.Op.Slug, "not-a-uuid")
	require.Error(t, err)

	_, _, err = lookupOperationEvidence(db, badOp.Op.Slug, goodOp.Evidence[0].UUID)
	require.Error(t, err)
}

func TestTagsForEvidenceByID(t *testing.T) {
	db := internalTestDBSetup(t)
	goodOp, _ := setupBasicTestOperation(t, db)

	// Add a tag for testing
	tagNames := []string{"one fish", "two fish", "red fish", "blue fish"}
	tagIDs := addTags(t, db, goodOp.ID, tagNames...)
	evidenceIDs := mapModelsEvidenceToID(goodOp.Evidence)
	require.True(t, len(tagIDs) >= len(evidenceIDs), "Future tests require that all tagIDs be used")
	zippedIDs := zipInt64Lists(tagIDs, evidenceIDs)
	linkTagToEvidence(t, db, zippedIDs)

	tagEvidenceMap, allTags, err := tagsForEvidenceByID(db, evidenceIDs)

	require.NoError(t, err)
	for i, ids := range zippedIDs {
		require.Equal(t, ids[0], tagEvidenceMap[ids[1]][0].ID)
		require.Equal(t, tagNames[i], tagEvidenceMap[ids[1]][0].Name)
		require.Equal(t, "black", tagEvidenceMap[ids[1]][0].ColorName)
	}
	// TODO: compare allTags to what was provided
	foundTagIDs := mapDtoTagToID(allTags)
	for _, id := range tagIDs {
		require.Contains(t, foundTagIDs, id)
	}

	require.NotNil(t, allTags)
}

func TestIsAdmin(t *testing.T) {
	adminCtx := middleware.InjectAdmin(context.Background(), true)
	nonAdminCtx := middleware.InjectAdmin(context.Background(), false)
	unknownCtx := context.Background()

	require.Nil(t, isAdmin(adminCtx))
	require.NotNil(t, isAdmin(nonAdminCtx))
	require.NotNil(t, isAdmin(unknownCtx))
}

func TestUserSlugToUserID(t *testing.T) {
	db := internalTestDBSetup(t)
	goodOp, _ := setupBasicTestOperation(t, db)

	userID, err := userSlugToUserID(db, goodOp.User.Slug)
	require.NoError(t, err)
	require.Equal(t, goodOp.UserID, userID)
}

func opDtoToModel(op *dtos.Operation, opID int64) models.Operation {
	return models.Operation{
		ID:     opID,
		Slug:   op.Slug,
		Name:   op.Name,
		Status: op.Status,
	}
}

func mapModelsEvidenceToID(in []models.Evidence) []int64 {
	rtn := make([]int64, len(in))
	for i := range in {
		rtn[i] = in[i].ID
	}
	return rtn
}

func mapDtoTagToID(in []dtos.Tag) []int64 {
	rtn := make([]int64, len(in))
	for i := range in {
		rtn[i] = in[i].ID
	}
	return rtn
}

func linkTagToEvidence(t *testing.T, db *database.Connection, tagIDEviIDZip [][2]int64) {
	err := db.BatchInsert("tag_evidence_map", len(tagIDEviIDZip), func(i int) map[string]interface{} {
		return map[string]interface{}{
			"tag_id":      tagIDEviIDZip[i][0],
			"evidence_id": tagIDEviIDZip[i][0],
		}
	})
	require.NoError(t, err)
}

func zipInt64Lists(a, b []int64) [][2]int64 {
	min := len(a)
	if len(b) < len(a) {
		min = len(b)
	}

	zip := make([][2]int64, min)
	for i := range zip {
		zip[i] = [2]int64{a[i], b[i]}
	}
	return zip
}

func addTags(t *testing.T, db *database.Connection, parentOpID int64, tagNames ...string) []int64 {
	err := db.BatchInsert("tags", len(tagNames), func(i int) map[string]interface{} {
		return map[string]interface{}{
			"operation_id": parentOpID,
			"color_name":   "black",
			"name":         tagNames[i],
		}
	})
	require.NoError(t, err)
	tagIDs := make([]int64, len(tagNames))
	for i, name := range tagNames {
		var id int64
		err := db.Get(&id, sq.Select("id").
			From("tags").
			Where(sq.Eq{"name": name, "operation_id": parentOpID}))
		require.NoError(t, err)
		tagIDs[i] = id
	}
	return tagIDs
}
