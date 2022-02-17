// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

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
