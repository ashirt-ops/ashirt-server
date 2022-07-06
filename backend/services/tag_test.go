// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestCreateTag(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)

	op := OpSorcerersStone
	i := services.CreateTagInput{
		Name:          "New Tag",
		ColorName:     "indigo",
		OperationSlug: op.Slug,
	}

	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	createdTag, err := services.CreateTag(ctx, db, i)
	require.NoError(t, err)
	require.Equal(t, createdTag.Name, i.Name)
	require.NotContains(t, HarryPotterSeedData.AllInitialTagIds(), createdTag.ID, "Should have new ID")

	updatedTag := getTagByID(t, db, createdTag.ID)

	require.Equal(t, op.ID, updatedTag.OperationID, "is in right operation")
}

func TestCreateDefaultTag(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	adminUser := UserDumbledore

	i := services.CreateDefaultTagInput{
		Name:      "New Tag",
		ColorName: "indigo",
	}

	// verify that a normal cannot create a new default tag
	ctx := simpleFullContext(normalUser)
	_, err := services.CreateDefaultTag(ctx, db, i)
	require.Error(t, err)

	// verify that an admin can create a new default tag
	ctx = simpleFullContext(adminUser)
	createdTag, err := services.CreateDefaultTag(ctx, db, i)
	require.NoError(t, err)
	require.Equal(t, createdTag.Name, i.Name)
	require.NotContains(t, HarryPotterSeedData.AllInitialDefaultTagIds(), createdTag.ID, "Should have new ID")
}

func TestDeleteTag(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)

	op := OpChamberOfSecrets
	i := services.DeleteTagInput{
		ID:            TagEarth.ID,
		OperationSlug: op.Slug,
	}

	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	err := services.DeleteTag(ctx, db, i)
	require.NoError(t, err)

	require.NotContains(t, getTagFromOperationID(t, db, op.ID), TagEarth, "TagEarth should have been deleted")
}

func TestDeleteDefaultTag(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	tagToRemove := DefaultTagWho
	normalUser := UserRon
	adminUser := UserDumbledore

	i := services.DeleteDefaultTagInput{
		ID: tagToRemove.ID,
	}

	// verify that a normal user cannot delete default tags
	ctx := simpleFullContext(normalUser)
	err := services.DeleteDefaultTag(ctx, db, i)
	require.Error(t, err)

	// verify that an admin can delete default tags
	ctx = simpleFullContext(adminUser)
	err = services.DeleteDefaultTag(ctx, db, i)
	require.NoError(t, err)
	require.NotContains(t, getDefaultTags(t, db), tagToRemove)
}

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

func TestListTagsByEvidenceDate(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	masterOp := OpGanttChart
	input := services.ListTagsByEvidenceDateInput{
		OperationSlug: masterOp.Slug,
	}

	expectedData := HarryPotterSeedData.TagIDsUsageByDate(masterOp.ID)

	// test no-access
	_, err := services.ListTagsByEvidenceDate(contextForUser(UserDraco, db), db, input)
	require.Error(t, err)

	// test read
	actualData, err := services.ListTagsByEvidenceDate(contextForUser(UserGinny, db), db, input)
	require.NoError(t, err)
	validateTagUsageData(t, HarryPotterSeedData, expectedData, actualData)

	// test write
	actualData, err = services.ListTagsByEvidenceDate(contextForUser(UserHarry, db), db, input)
	require.NoError(t, err)
	validateTagUsageData(t, HarryPotterSeedData, expectedData, actualData)

	// test admin
	actualData, err = services.ListTagsByEvidenceDate(contextForUser(UserDumbledore, db), db, input)
	require.NoError(t, err)
	validateTagUsageData(t, HarryPotterSeedData, expectedData, actualData)
}

func TestListTagsForOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	allTags := getTagFromOperationID(t, db, masterOp.ID)
	require.NotEqual(t, len(allTags), 0, "Some number of tags should exist")

	tags, err := services.ListTagsForOperation(ctx, db, services.ListTagsForOperationInput{masterOp.Slug})
	require.NoError(t, err)
	require.Equal(t, len(tags), len(allTags))

	dtoTags := make([]*dtos.Tag, len(tags))
	for i, tag := range tags {
		dtoTags[i] = &tag.Tag
		require.Equal(t, tag.EvidenceCount, getTagUsage(t, db, tag.ID))
	}

	validateTagSets(t, dtoTags, allTags, validateTag)
}

func validateTag(t *testing.T, expected models.Tag, actual *dtos.Tag) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.ColorName, actual.ColorName)
}

func validateTagSets(t *testing.T, dtoSet []*dtos.Tag, dbSet []models.Tag, validate func(*testing.T, models.Tag, *dtos.Tag)) {
	var expected *models.Tag = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.ID == dtoItem.ID {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validate(t, *expected, dtoItem)
	}
}

func ptrTagListToReal(in []*dtos.Tag) []dtos.Tag {
	rtn := make([]dtos.Tag, len(in))
	for i, item := range in {
		rtn[i] = *item
	}
	return rtn
}

func realTagListToPtr(in []dtos.Tag) []*dtos.Tag {
	rtn := make([]*dtos.Tag, len(in))
	for i, item := range in {
		rtn[i] = &item
	}
	return rtn
}

func TestListDefaultTags(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	normalUser := UserRon
	adminUser := UserDumbledore

	allTags := getDefaultTags(t, db)
	require.NotEqual(t, len(allTags), 0, "Some number of default tags should exist")

	// verify that normal users cannot list default tags
	ctx := simpleFullContext(normalUser)
	_, err := services.ListDefaultTags(ctx, db)
	require.Error(t, err)

	// verify that admins can list default tags
	ctx = simpleFullContext(adminUser)
	tags, err := services.ListDefaultTags(ctx, db)
	require.NoError(t, err)
	require.Equal(t, len(tags), len(allTags))

	validateDefaultTagSets(t, tags, allTags, validateDefaultTag)
}

func validateDefaultTagSets(
	t *testing.T,
	dtoSet []*dtos.DefaultTag,
	dbSet []models.DefaultTag,
	validate func(*testing.T, models.DefaultTag, *dtos.DefaultTag),
) {
	var expected *models.DefaultTag = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.ID == dtoItem.ID {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validate(t, *expected, dtoItem)
	}
}

func validateDefaultTag(t *testing.T, expected models.DefaultTag, actual *dtos.DefaultTag) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.ColorName, actual.ColorName)
}

func validateTagUsageData(t *testing.T, seed TestSeedData, expectedTagIDUsage map[int64][]time.Time, actual []*dtos.TagByEvidenceDate) {
	require.Equal(t, len(expectedTagIDUsage), len(actual))
	for tagID, dates := range expectedTagIDUsage {
		tag := seed.GetTagFromID(tagID)

		var match *dtos.TagByEvidenceDate = nil
		for _, tagUsage := range actual {
			if tagUsage.ID == tag.ID {
				match = tagUsage
				break
			}
		}
		require.NotNil(t, match)
		require.Equal(t, tag.ColorName, match.ColorName)
		require.Equal(t, tag.Name, match.Name)

		mapDates := func(times []time.Time) map[int64]int {
			rtn := make(map[int64]int)
			for _, someTime := range times {
				timeCount := 0
				for _, matchTime := range times {
					if someTime == matchTime {
						timeCount += 1
					}
				}
				rtn[someTime.Unix()] = timeCount
			}
			return rtn
		}

		expectedDateUsage := mapDates(dates)
		actualDateUsage := mapDates(match.UsageDates)
		require.Equal(t, expectedDateUsage, actualDateUsage)
	}
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

func TestUpdateTag(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	op := OpChamberOfSecrets
	i := services.UpdateTagInput{
		ID:            TagEarth.ID,
		OperationSlug: op.Slug,
		Name:          "Moon",
		ColorName:     "green",
	}

	ctx := fullContext(UserHarry.ID, &policy.FullAccess{})
	err := services.UpdateTag(ctx, db, i)
	require.NoError(t, err)

	updatedTag := getTagByID(t, db, TagEarth.ID)
	require.Equal(t, models.Tag{
		ID:          TagEarth.ID,
		OperationID: op.ID,
		Name:        "Moon",
		ColorName:   "green",
		CreatedAt:   TagEarth.CreatedAt,
		UpdatedAt:   updatedTag.UpdatedAt,
	}, updatedTag)
}

func TestUpdateDefaultTag(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	normalUser := UserRon
	adminUser := UserDumbledore
	tagToUpdate := DefaultTagWho

	i := services.UpdateDefaultTagInput{
		ID:        tagToUpdate.ID,
		Name:      "How",
		ColorName: "green",
	}

	// verify that a normal user cannot update a default tags
	ctx := simpleFullContext(normalUser)
	err := services.UpdateDefaultTag(ctx, db, i)
	require.Error(t, err)

	// verify that an admin can update default tags
	ctx = simpleFullContext(adminUser)
	err = services.UpdateDefaultTag(ctx, db, i)
	require.NoError(t, err)

	updatedTag := getDefaultTagByID(t, db, tagToUpdate.ID)
	require.Equal(t, models.DefaultTag{
		ID:        tagToUpdate.ID,
		Name:      i.Name,
		ColorName: i.ColorName,
		CreatedAt: tagToUpdate.CreatedAt,
		UpdatedAt: updatedTag.UpdatedAt,
	}, updatedTag)
}
