// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"

	sq "github.com/Masterminds/squirrel"
)

func TestReadEvidence(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})
	cs, _ := contentstore.NewMemStore()

	masterOp := OpChamberOfSecrets
	masterEvidence := cloneEvidence(EviFlyingCar) // cloning so we don't alter the underlying evidence for other tests

	neitherInput := services.ReadEvidenceInput{
		OperationSlug: masterOp.Slug,
		EvidenceUUID:  masterEvidence.UUID,
	}

	retrievedEvidence, err := services.ReadEvidence(ctx, db, cs, neitherInput)
	require.NoError(t, err)
	validateReadEvidenceOutput(t, masterEvidence, retrievedEvidence)

	fullImg := []byte("full_image")
	thumbImg := []byte("thumb_image")
	thumbKey, err := cs.Upload(bytes.NewReader(thumbImg))
	require.NoError(t, err)
	fullKey, err := cs.Upload(bytes.NewReader(fullImg))
	require.NoError(t, err)

	err = db.Update(sq.Update("evidence").
		SetMap(map[string]interface{}{
			"full_image_key":  fullKey,
			"thumb_image_key": thumbKey,
		}).
		Where(sq.Eq{"id": masterEvidence.ID}))
	require.NoError(t, err)
	masterEvidence.FullImageKey = fullKey
	masterEvidence.ThumbImageKey = thumbKey

	retrievedEvidence, err = services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{
		OperationSlug: masterOp.Slug,
		EvidenceUUID:  masterEvidence.UUID,
		LoadPreview:   true,
	})
	require.NoError(t, err)
	validateReadEvidenceOutput(t, masterEvidence, retrievedEvidence)
	previewBytes, err := ioutil.ReadAll(retrievedEvidence.Preview)
	require.NoError(t, err)
	require.Equal(t, thumbImg, previewBytes)

	retrievedEvidence, err = services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{
		OperationSlug: masterOp.Slug,
		EvidenceUUID:  masterEvidence.UUID,
		LoadMedia:     true,
	})
	require.NoError(t, err)
	validateReadEvidenceOutput(t, masterEvidence, retrievedEvidence)
	mediaBytes, err := ioutil.ReadAll(retrievedEvidence.Media)
	require.NoError(t, err)
	require.Equal(t, fullImg, mediaBytes)
}

func cloneEvidence(in models.Evidence) models.Evidence {
	return models.Evidence{
		ID:            in.ID,
		UUID:          in.UUID,
		OperationID:   in.OperationID,
		OperatorID:    in.OperatorID,
		Description:   in.Description,
		ContentType:   in.ContentType,
		FullImageKey:  in.FullImageKey,
		ThumbImageKey: in.ThumbImageKey,
		OccurredAt:    in.OccurredAt,
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
	}
}

func validateReadEvidenceOutput(t *testing.T, expected models.Evidence, actual *services.ReadEvidenceOutput) {
	require.Equal(t, expected.UUID, actual.UUID)
	require.Equal(t, expected.Description, actual.Description)
	require.Equal(t, expected.ContentType, actual.ContentType)
	require.Equal(t, expected.OccurredAt, actual.OccurredAt)
}
