// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"bytes"
	// "encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateEvidence(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	memStore, _ := contentstore.NewMemStore()
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

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
	validateInsertedEvidence(t, imgEvi, imgInput, err, op, memStore, db, imgContent)

	cbContent := []byte("I'm a codeblock!")
	cbInput := services.CreateEvidenceInput{
		OperationSlug: op.Slug,
		Description:   "some codeblock",
		ContentType:   "codeblock",
		TagIDs:        TagIDsFromTags(TagVenus, TagMars),
		Content:       bytes.NewReader(cbContent),
	}
	cbEvi, err := services.CreateEvidence(ctx, db, memStore, cbInput)
	validateInsertedEvidence(t, cbEvi, cbInput, err, op, memStore, db, cbContent)

	bareInput := services.CreateEvidenceInput{
		OperationSlug: op.Slug,
		Description:   "Just a note here",
		ContentType:   "Plain Text",
		TagIDs:        TagIDsFromTags(TagVenus, TagMars),
	}
	bareEvi, err := services.CreateEvidence(ctx, db, memStore, bareInput)
	validateInsertedEvidence(t, bareEvi, bareInput, err, op, memStore, db, nil)
}

func validateInsertedEvidence(t *testing.T, evi *dtos.Evidence, src services.CreateEvidenceInput,
	err error, op models.Operation, store *contentstore.MemStore, db *database.Connection,
	rawContent []byte) {

	require.NoError(t, err)
	fullEvidence := getEvidenceByUUID(t, db, evi.UUID)
	require.Equal(t, src.Description, fullEvidence.Description)
	require.Equal(t, src.ContentType, fullEvidence.ContentType)
	require.Equal(t, UserRon.ID, fullEvidence.OperatorID)
	require.Equal(t, op.ID, fullEvidence.OperationID, "Associated with right operation")
	if src.Content != nil {
		require.NotEqual(t, "", fullEvidence.FullImageKey, "Keys should be populated")
		require.NotEqual(t, "", fullEvidence.ThumbImageKey, "Keys should be populated")
		fullReader, _ := store.Read(fullEvidence.FullImageKey)
		fullContentBytes, _ := ioutil.ReadAll(fullReader)
		require.Equal(t, rawContent, fullContentBytes)
	}

	tagIDs := getTagIDsFromEvidenceID(t, db, fullEvidence.ID)

	require.Equal(t, sorted(tagIDs), sorted(src.TagIDs))
}
