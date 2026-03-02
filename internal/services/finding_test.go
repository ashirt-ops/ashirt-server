package services_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/dtos"
	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/services"
	"github.com/stretchr/testify/require"
)

type findingValidator func(*testing.T, models.Finding, *dtos.Finding)

func TestCreateFinding(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		op := OpChamberOfSecrets
		i := services.CreateFindingInput{
			OperationSlug: op.Slug,
			Category:      VendorFindingCategory.Category,
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
	})
}

func TestDeleteFinding(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

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
	})
}

func TestListFindingsForOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpChamberOfSecrets
		input := services.ListFindingsForOperationInput{
			OperationSlug: masterOp.Slug,
			Filters:       helpers.TimelineFilters{},
		}

		allFindings := getFindingsByOperationID(t, db, masterOp.ID)
		require.NotEqual(t, len(allFindings), 0, "Some number of findings should exist")

		foundFindings, err := services.ListFindingsForOperation(ctx, db, input)
		require.NoError(t, err)
		require.Equal(t, len(foundFindings), len(allFindings))
		validateFindingSets(t, foundFindings, allFindings, buildFindingValidator(seed))
	})
}

func TestAddEvidenceToFinding(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpChamberOfSecrets
		masterFinding := FindingBook2Magic
		evidenceToAdd1 := EviSpiderAragog
		evidenceToAdd2 := EviMoaningMyrtle
		evidenceToRemove1 := EviDobby
		evidenceToRemove2 := EviFlyingCar

		initialEvidenceList := getEvidenceIDsFromFinding(t, db, masterFinding.ID)

		expectedEvidenceSet := make(map[int64]bool)
		for _, id := range initialEvidenceList {
			if id != evidenceToRemove1.ID && id != evidenceToRemove2.ID {
				expectedEvidenceSet[id] = true
			}
		}
		expectedEvidenceSet[evidenceToAdd1.ID] = true
		expectedEvidenceSet[evidenceToAdd2.ID] = true
		expectedEvidenceList := make([]int64, 0, len(expectedEvidenceSet))
		for key, v := range expectedEvidenceSet {
			if v {
				expectedEvidenceList = append(expectedEvidenceList, key)
			}
		}

		i := services.AddEvidenceToFindingInput{
			OperationSlug:    masterOp.Slug,
			FindingUUID:      masterFinding.UUID,
			EvidenceToAdd:    []string{evidenceToAdd1.UUID, evidenceToAdd2.UUID},
			EvidenceToRemove: []string{evidenceToRemove1.UUID, evidenceToRemove2.UUID},
		}
		err := services.AddEvidenceToFinding(ctx, db, i)
		require.NoError(t, err)

		changedEvidenceList := getEvidenceIDsFromFinding(t, db, masterFinding.ID)

		require.Equal(t, sorted(expectedEvidenceList), sorted(changedEvidenceList))
	})
}

func TestReadFinding(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserRon, db)

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
		require.Equal(t, seed.CategoryForFinding(masterFinding), retrievedFinding.Category)
		require.Equal(t, masterFinding.Description, retrievedFinding.Description)
		require.Equal(t, masterFinding.ReadyToReport, retrievedFinding.ReadyToReport)
		require.Equal(t, masterFinding.TicketLink, retrievedFinding.TicketLink)
		require.Equal(t, len(seed.EvidenceIDsForFinding(masterFinding)), retrievedFinding.NumEvidence)
		validateTagSets(t, realTagListToPtr(retrievedFinding.Tags), seed.TagsForFinding(masterFinding), validateTag)
	})
}

func realTagListToPtr(in []dtos.Tag) []*dtos.Tag {
	return helpers.Map(in, func(t dtos.Tag) *dtos.Tag { return &t })
}

func TestUpdateFinding(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		// tests for common fields
		masterOp := OpChamberOfSecrets
		masterFinding := FindingBook2Magic
		input := services.UpdateFindingInput{
			OperationSlug: masterOp.Slug,
			FindingUUID:   masterFinding.UUID,
			Category:      DetectionGapFindingCategory.Category,
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
	})
}

func buildFindingValidator(seed TestSeedData) findingValidator {
	return func(t *testing.T, expected models.Finding, actual *dtos.Finding) {
		require.Equal(t, expected.UUID, actual.UUID)
		require.Equal(t, seed.CategoryForFinding(expected), actual.Category)
		require.Equal(t, expected.Title, actual.Title)
		require.Equal(t, expected.Description, actual.Description)
		require.Equal(t, expected.ReadyToReport, actual.ReadyToReport)
		require.Equal(t, expected.TicketLink, actual.TicketLink)
	}
}

func validateFindingSets(t *testing.T, dtoSet []*dtos.Finding, dbSet []models.Finding, validate findingValidator) {
	var expected *models.Finding = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.UUID == dtoItem.UUID {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validate(t, *expected, dtoItem)
	}
}
