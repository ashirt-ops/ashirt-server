// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateEvidence(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)

	memStore, _ := contentstore.NewMemStore()
	normalUser := UserRon
	ctx := contextForUser(normalUser, db)

	op := OpChamberOfSecrets
	imgContent := TinyImg
	imgInput := services.CreateEvidenceInput{
		OperationSlug: op.Slug,
		Description:   "some image",
		ContentType:   "image",
		TagIDs:        TagIDsFromTags(TagSaturn, TagJupiter),
		Content:       bytes.NewReader(imgContent),
	}
	imgEvi, err := services.CreateEvidence(ctx, db, memStore, imgInput)
	require.NoError(t, err)
	validateInsertedEvidence(t, imgEvi, imgInput, normalUser, op, memStore, db, imgContent)

	cbContent := []byte("I'm a codeblock!")
	cbInput := services.CreateEvidenceInput{
		OperationSlug: op.Slug,
		Description:   "some codeblock",
		ContentType:   "codeblock",
		TagIDs:        TagIDsFromTags(TagVenus, TagMars),
		Content:       bytes.NewReader(cbContent),
	}
	cbEvi, err := services.CreateEvidence(ctx, db, memStore, cbInput)
	require.NoError(t, err)
	validateInsertedEvidence(t, cbEvi, cbInput, normalUser, op, memStore, db, cbContent)

	bareInput := services.CreateEvidenceInput{
		OperationSlug: op.Slug,
		Description:   "Just a note here",
		ContentType:   "Plain Text",
		TagIDs:        TagIDsFromTags(TagVenus, TagMars),
	}
	bareEvi, err := services.CreateEvidence(ctx, db, memStore, bareInput)
	require.NoError(t, err)
	validateInsertedEvidence(t, bareEvi, bareInput, normalUser, op, memStore, db, nil)
}

func TestHeadlessUserAccess(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	memStore, _ := contentstore.NewMemStore()

	headlessUser := UserHeadlessNick
	op := OpChamberOfSecrets

	// Pre-test: verify that the user is not a normal member of the operation
	opRoles := getUserRolesForOperationByOperationID(t, db, op.ID)
	for _, v := range opRoles {
		require.NotEqual(t, v.UserID, headlessUser.ID, "Pretest error: User should not have normal access to this operation")
	}

	ctx := contextForUser(headlessUser, db)

	imgContent := TinyImg
	imgInput := services.CreateEvidenceInput{
		OperationSlug: op.Slug,
		Description:   "headless image",
		ContentType:   "image",
		TagIDs:        TagIDsFromTags(TagSaturn, TagJupiter),
		Content:       bytes.NewReader(imgContent),
	}
	imgEvi, err := services.CreateEvidence(ctx, db, memStore, imgInput)
	require.NoError(t, err)
	fullEvidence := getEvidenceByUUID(t, db, imgEvi.UUID)
	require.Equal(t, imgInput.Description, fullEvidence.Description)
}

func validateInsertedEvidence(t *testing.T, evi *dtos.Evidence, src services.CreateEvidenceInput,
	user models.User, op models.Operation, store *contentstore.MemStore, db *database.Connection,
	rawContent []byte) {

	fullEvidence := getEvidenceByUUID(t, db, evi.UUID)
	require.Equal(t, src.Description, fullEvidence.Description)
	require.Equal(t, src.ContentType, fullEvidence.ContentType)
	require.Equal(t, user.ID, fullEvidence.OperatorID)
	require.Equal(t, op.ID, fullEvidence.OperationID, "Associated with right operation")
	if src.Content != nil {
		require.NotEqual(t, "", fullEvidence.FullImageKey, "Keys should be populated")
		require.NotEqual(t, "", fullEvidence.ThumbImageKey, "Keys should be populated")
		fullReader, _ := store.Read(fullEvidence.FullImageKey)
		fullContentBytes, _ := io.ReadAll(fullReader)
		require.Equal(t, rawContent, fullContentBytes)
	}

	tagIDs := getTagIDsFromEvidenceID(t, db, fullEvidence.ID)

	require.Equal(t, sorted(tagIDs), sorted(src.TagIDs))
}
