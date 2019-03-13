// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestListEvidenceForFinding(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	masterFinding := FindingBook2Magic
	allEvidence := getFullEvidenceByFindingID(t, db, masterFinding.ID)

	require.NotEqual(t, 0, len(allEvidence), "Some evidence should be present for this finding")

	input := services.ListEvidenceForFindingInput{
		OperationSlug: masterOp.Slug,
		FindingUUID:   FindingBook2Magic.UUID,
	}

	foundEvidence, err := services.ListEvidenceForFinding(ctx, db, input)
	require.NoError(t, err)
	require.Equal(t, len(foundEvidence), len(allEvidence))
	validateEvidenceSets(t, foundEvidence, allEvidence, validateEvidence)
}

type evidenceValidator func(*testing.T, FullEvidence, dtos.Evidence)

func validateEvidence(t *testing.T, expected FullEvidence, actual dtos.Evidence) {
	require.Equal(t, expected.UUID, actual.UUID)
	require.Equal(t, expected.ContentType, actual.ContentType)
	require.Equal(t, expected.Description, actual.Description)
	validateTagSets(t, toPtrTagList(actual.Tags), expected.Tags, validateTag)
	require.Equal(t, expected.OccurredAt, actual.OccurredAt)

	require.Equal(t, expected.Slug, actual.Operator.Slug)
	require.Equal(t, expected.FirstName, actual.Operator.FirstName)
	require.Equal(t, expected.LastName, actual.Operator.LastName)
}

func validateEvidenceSets(t *testing.T, dtoSet []dtos.Evidence, dbSet []FullEvidence, validator evidenceValidator) {
	var expected *FullEvidence = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.UUID == dtoItem.UUID {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validator(t, *expected, dtoItem)
	}
}

func toPtrTagList(in []dtos.Tag) []*dtos.Tag {
	rtn := make([]*dtos.Tag, len(in))

	for i := range in {
		rtn[i] = &in[i]
	}

	return rtn
}
