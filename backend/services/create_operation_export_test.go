// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/services"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func TestCreateOperationExport(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := simpleFullContext(UserRon)

	targetOperation := OpChamberOfSecrets

	// verify feature locked to super admins only
	output := services.CreateOperationExport(ctx, db, targetOperation.Slug)
	require.Error(t, output.Err)

	ctx = simpleFullContext(UserDumbledore)

	// This area is safe from the export processor because it does act on the normal database,
	// So the processor should not find these records
	output = services.CreateOperationExport(ctx, db, targetOperation.Slug)
	require.NoError(t, output.Err)
	require.True(t, output.Queued)

	// repeat with the same values to ensure that items are not re-queued while waiting to be
	// processed.
	output = services.CreateOperationExport(ctx, db, targetOperation.Slug)
	require.NoError(t, output.Err)
	require.False(t, output.Queued)

	// repeat once more, after setting status to something besides "Pending"
	db.Update(sq.Update("exports_queue").Set("status", models.ExportStatusComplete)) //acting on all, because there should only be 1
	output = services.CreateOperationExport(ctx, db, targetOperation.Slug)
	require.NoError(t, output.Err)
	require.True(t, output.Queued)
}
