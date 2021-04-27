// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestReadFinding(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets
	masterFinding := FindingBook2Magic

	input := services.ReadFindingInput{
		OperationSlug: masterOp.Slug,
		FindingUUID:   masterFinding.UUID,
	}

	retrievedFinding, err := services.ReadFinding(ctx, db, input)
	require.NoError(t, err)

	require.Equal(t, masterFinding.UUID, retrievedFinding.UUID)
	require.Equal(t, masterFinding.Title, retrievedFinding.Title)
	require.Equal(t, HarryPotterSeedData.CategoryForFinding(masterFinding), retrievedFinding.Category)
	require.Equal(t, masterFinding.Description, retrievedFinding.Description)
	require.Equal(t, masterFinding.ReadyToReport, retrievedFinding.ReadyToReport)
	require.Equal(t, masterFinding.TicketLink, retrievedFinding.TicketLink)
	require.Equal(t, len(HarryPotterSeedData.EvidenceIDsForFinding(masterFinding)), retrievedFinding.NumEvidence)
	validateTagSets(t, realTagListToPtr(retrievedFinding.Tags), HarryPotterSeedData.TagsForFinding(masterFinding), validateTag)
}
