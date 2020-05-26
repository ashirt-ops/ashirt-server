// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
)

func TestDeleteFinding(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	finding := FindingBook2Magic
	i := services.DeleteFindingInput{
		OperationSlug: OpChamberOfSecrets.Slug,
		FindingUUID:   finding.UUID,
	}
	getFindingCount := makeDBRowCounter(t, db, "findings", "id=?", finding.ID)
	require.Equal(t, int64(1), getFindingCount(), "Database should have a finding to delete")

	getMappedEvidenceCount := makeDBRowCounter(t, db, "evidence_finding_map", "finding_id=?", finding.ID)
	require.True(t, getMappedEvidenceCount() > 0, "Database should have associated evidence to delete")

	err := services.DeleteFinding(ctx, db, i)
	require.NoError(t, err)
	require.Equal(t, int64(0), getFindingCount(), "Database should have deleted the finding")
	require.Equal(t, int64(0), getMappedEvidenceCount(), "Database should have deleted evidence mapping")
}
