// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestCreateFinding(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	op := OpChamberOfSecrets
	i := services.CreateFindingInput{
		OperationSlug: op.Slug,
		Category:      "Garbage",
		Title:         "When Dinosaurs Attack",
		Description:   "An investigative look into what happens when dinosaurs vandalize neighborhoods like yours",
	}
	createdFinding, err := services.CreateFinding(ctx, db, i)
	require.NoError(t, err)
	fullFinding, err := services.ReadFinding(ctx, db, services.ReadFindingInput{OperationSlug: op.Slug, FindingUUID: createdFinding.UUID})
	require.NoError(t, err)

	require.Equal(t, i.Category, fullFinding.Category)
	require.Equal(t, i.Title, fullFinding.Title)
	require.Equal(t, i.Description, fullFinding.Description)
}
