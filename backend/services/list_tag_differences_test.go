// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestListTagDifferences(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	startingOp := OpChamberOfSecrets
	endingOp := OpSorcerersStone

	input := services.ListTagsDifferenceInput{
		SourceOperationSlug:      startingOp.Slug,
		DestinationOperationSlug: endingOp.Slug,
	}

	// verify that Neville (cannot read endingOp) cannot determine tag differences
	ctx := contextForUser(UserNeville, db)
	_, err := services.ListTagDifference(ctx, db, input)
	require.Error(t, err)

	// verify that Harry (can read both) can determine tag differences
	ctx = contextForUser(UserHarry, db)
	data, err := services.ListTagDifference(ctx, db, input)
	require.NoError(t, err)

	verifySharedTags(t, data, CommonTagWhoCoS, CommonTagWhatCoS, CommonTagWhereCoS, CommonTagWhenCoS, CommonTagWhyCoS)
	verifyDroppedTags(t, data, TagMercury, TagVenus, TagEarth, TagMars, TagJupiter, TagSaturn, TagNeptune)
}

func TestListTagDifferencesForEvidence(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	startingOp := OpChamberOfSecrets
	endingOp := OpSorcerersStone
	sourceEvidence := EviPetrifiedHermione //shares tags between the two operations

	input := services.ListTagDifferenceForEvidenceInput{
		ListTagsDifferenceInput: services.ListTagsDifferenceInput{
			SourceOperationSlug:      startingOp.Slug,
			DestinationOperationSlug: endingOp.Slug,
		},
		SourceEvidenceUUID: sourceEvidence.UUID,
	}

	// verify that Neville (cannot read endingOp) cannot determine tag differences
	ctx := contextForUser(UserNeville, db)
	_, err := services.ListTagDifferenceForEvidence(ctx, db, input)
	require.Error(t, err)

	// verify that Harry (can read both) can determine tag differences
	ctx = contextForUser(UserHarry, db)
	data, err := services.ListTagDifferenceForEvidence(ctx, db, input)
	require.NoError(t, err)

	verifySharedTags(t, data, CommonTagWhatCoS, CommonTagWhoCoS)
	verifyDroppedTags(t, data, TagMars)
}

func verifySharedTags(t *testing.T, diff *dtos.TagDifference, sharedTags ...models.Tag) {
	foundTags := make([]bool, len(sharedTags))
	extraTags := make([]int64, 0, len(diff.Included))

	for _, tagpair := range diff.Included {
		foundMatch := false
		for tagIndex, modelTag := range sharedTags {
			if tagpair.SourceTag.ID == modelTag.ID {
				foundTags[tagIndex] = true
				foundMatch = true
			}
		}
		if !foundMatch {
			extraTags = append(extraTags, tagpair.SourceTag.ID)
		}
	}
	require.True(t, len(extraTags) == 0, "Ensure no extra tags are present")
	allFound := true
	for _, v := range foundTags {
		allFound = allFound && v
	}
	require.True(t, allFound)
}

func verifyDroppedTags(t *testing.T, diff *dtos.TagDifference, separateTags ...models.Tag) {
	foundTags := make([]bool, len(separateTags))
	extraTags := make([]int64, 0, len(diff.Included))

	for _, diffedTag := range diff.Excluded {
		foundMatch := false
		for tagIndex, modelTag := range separateTags {
			if diffedTag.ID == modelTag.ID {
				foundTags[tagIndex] = true
				foundMatch = true
			}
		}
		if !foundMatch {
			extraTags = append(extraTags, diffedTag.ID)
		}
	}
	require.True(t, len(extraTags) == 0, "Ensure no extra tags are present")
	allFound := true
	for _, v := range foundTags {
		allFound = allFound && v
	}
	require.True(t, allFound)
}
