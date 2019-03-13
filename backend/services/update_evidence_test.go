// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestUpdateEvidence(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})
	cs, _ := contentstore.NewMemStore()

	// tests for common fields
	masterOp := OpChamberOfSecrets
	masterEvidence := EviFlyingCar
	initialTags := HarryPotterSeedData.TagsForEvidence(masterEvidence)
	tagToAdd := TagMercury
	tagToRemove := TagSaturn
	description := "New Description"
	input := services.UpdateEvidenceInput{
		OperationSlug: masterOp.Slug,
		EvidenceUUID:  masterEvidence.UUID,
		Description:   &description,
		TagsToRemove:  []int64{tagToRemove.ID},
		TagsToAdd:     []int64{tagToAdd.ID},
	}
	require.Contains(t, initialTags, tagToRemove)
	require.NotContains(t, initialTags, tagToAdd)

	err := services.UpdateEvidence(ctx, db, cs, input)
	require.NoError(t, err)
	evi, err := services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{OperationSlug: masterOp.Slug, EvidenceUUID: masterEvidence.UUID})
	require.NoError(t, err)
	require.Equal(t, *input.Description, evi.Description)
	expectedTagIDs := make([]int64, 0, len(initialTags))
	for _, t := range initialTags {
		if t != tagToRemove {
			expectedTagIDs = append(expectedTagIDs, t.ID)
		}
	}
	expectedTagIDs = append(expectedTagIDs, tagToAdd.ID)
	require.Equal(t, sorted(expectedTagIDs), sorted(getTagIDsFromEvidenceID(t, db, masterEvidence.ID)))

	// test for content

	codeblockEvidence := EviTomRiddlesDiary
	newContent := "stabbed_with_basilisk_fang = False\n\ndef is_alive():\n  return not stabbed_with_basilisk_fang\n"
	input = services.UpdateEvidenceInput{
		OperationSlug: masterOp.Slug,
		EvidenceUUID:  codeblockEvidence.UUID,
		Description:   &codeblockEvidence.Description, // Note: A quirk with UpdateEvidence is that it will always update the description, even if it is empty.
		Content:       bytes.NewReader([]byte(newContent)),
	}

	err = services.UpdateEvidence(ctx, db, cs, input)
	require.NoError(t, err)
	evi, err = services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{
		OperationSlug: masterOp.Slug,
		EvidenceUUID:  codeblockEvidence.UUID,
		LoadMedia:     true,
		LoadPreview:   true,
	})
	require.NoError(t, err)
	mediaBytes, err := ioutil.ReadAll(evi.Media)
	require.NoError(t, err)
	previewBytes, err := ioutil.ReadAll(evi.Preview)
	require.NoError(t, err)
	require.Equal(t, mediaBytes, previewBytes, "Preview and Media content should be identical for codeblocks")
	require.Equal(t, []byte(newContent), previewBytes)

	updatedEvidence := getEvidenceByID(t, db, codeblockEvidence.ID)
	require.Equal(t, updatedEvidence.ThumbImageKey, updatedEvidence.FullImageKey)
	require.NotEqual(t, "", updatedEvidence.FullImageKey)
}
