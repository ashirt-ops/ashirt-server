// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/helpers"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestListEvidenceForOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	allEvidence := getFullEvidenceByOperationID(t, db, masterOp.ID)

	require.NotEqual(t, len(allEvidence), 0, "Some evidence should be present")

	input := services.ListEvidenceForOperationInput{
		OperationSlug: masterOp.Slug,
		Filters:       helpers.TimelineFilters{},
	}

	foundEvidence, err := services.ListEvidenceForOperation(ctx, db, input)
	require.NoError(t, err)
	require.Equal(t, len(foundEvidence), len(allEvidence))
	validateEvidenceSets(t, toRealEvidenceList(foundEvidence), allEvidence, validateEvidence)
}

func toRealEvidenceList(in []*dtos.Evidence) []dtos.Evidence {
	rtn := make([]dtos.Evidence, len(in))

	for i := range in {
		rtn[i] = *in[i]
	}

	return rtn
}
