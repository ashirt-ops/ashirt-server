// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestUpdateFinding(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	// tests for common fields
	masterOp := OpChamberOfSecrets
	masterFinding := FindingBook2Magic
	input := services.UpdateFindingInput{
		OperationSlug: masterOp.Slug,
		FindingUUID:   masterFinding.UUID,
		Category:      "New Category",
		Title:         "New Title",
		Description:   "New Description",
	}

	err := services.UpdateFinding(ctx, db, input)
	require.NoError(t, err)
	finding, err := services.ReadFinding(ctx, db, services.ReadFindingInput{OperationSlug: masterOp.Slug, FindingUUID: masterFinding.UUID})
	require.NoError(t, err)
	require.Equal(t, input.Description, finding.Description)
	require.Equal(t, input.Title, finding.Title)
	require.Equal(t, input.Category, finding.Category)
}
