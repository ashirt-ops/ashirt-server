package services_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

type evidenceValidator func(*testing.T, FullEvidence, dtos.Evidence)

func TestCreateEvidence(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		memStore, _ := contentstore.NewMemStore()
		normalUser := UserRon
		ctx := contextForUser(normalUser, db)

		op := OpChamberOfSecrets
		imgContent := TinyImg
		imgInput := services.CreateEvidenceInput{
			OperationSlug: op.Slug,
			Description:   "some image",
			ContentType:   "image",
			TagIDs:        TagIDsFromTags(TagSaturn, TagJupiter),
			Content:       bytes.NewReader(imgContent),
		}
		imgEvi, err := services.CreateEvidence(ctx, db, memStore, imgInput)
		require.NoError(t, err)
		validateInsertedEvidence(t, imgEvi, imgInput, normalUser, op, memStore, db, imgContent)

		cbContent := []byte("I'm a codeblock!")
		cbInput := services.CreateEvidenceInput{
			OperationSlug: op.Slug,
			Description:   "some codeblock",
			ContentType:   "codeblock",
			TagIDs:        TagIDsFromTags(TagVenus, TagMars),
			Content:       bytes.NewReader(cbContent),
		}
		cbEvi, err := services.CreateEvidence(ctx, db, memStore, cbInput)
		require.NoError(t, err)
		validateInsertedEvidence(t, cbEvi, cbInput, normalUser, op, memStore, db, cbContent)

		bareInput := services.CreateEvidenceInput{
			OperationSlug: op.Slug,
			Description:   "Just a note here",
			ContentType:   "Plain Text",
			TagIDs:        TagIDsFromTags(TagVenus, TagMars),
		}
		bareEvi, err := services.CreateEvidence(ctx, db, memStore, bareInput)
		require.NoError(t, err)
		validateInsertedEvidence(t, bareEvi, bareInput, normalUser, op, memStore, db, nil)
	})
}

func TestHeadlessUserAccess(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		memStore, _ := contentstore.NewMemStore()

		headlessUser := UserHeadlessNick
		op := OpChamberOfSecrets

		// Pre-test: verify that the user is not a normal member of the operation
		opRoles := getUserRolesForOperationByOperationID(t, db, op.ID)
		for _, v := range opRoles {
			require.NotEqual(t, v.UserID, headlessUser.ID, "Pretest error: User should not have normal access to this operation")
		}

		ctx := contextForUser(headlessUser, db)

		imgContent := TinyImg
		imgInput := services.CreateEvidenceInput{
			OperationSlug: op.Slug,
			Description:   "headless image",
			ContentType:   "image",
			TagIDs:        TagIDsFromTags(TagSaturn, TagJupiter),
			Content:       bytes.NewReader(imgContent),
		}
		imgEvi, err := services.CreateEvidence(ctx, db, memStore, imgInput)
		require.NoError(t, err)
		fullEvidence := getEvidenceByUUID(t, db, imgEvi.UUID)
		require.Equal(t, imgInput.Description, fullEvidence.Description)
	})
}

func TestDeleteEvidenceNoPropogate(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		op := OpChamberOfSecrets
		memStore := createPopulatedMemStore(seed)

		masterEvidence := EviFlyingCar
		i := services.DeleteEvidenceInput{
			OperationSlug:            op.Slug,
			EvidenceUUID:             masterEvidence.UUID,
			DeleteAssociatedFindings: false,
		}
		// populate content store
		contentStoreKey := masterEvidence.UUID // seed data shares full and thumb key ids

		getAssociatedTagCount := makeDBRowCounter(t, db, "tag_evidence_map", "evidence_id=?", masterEvidence.ID)
		require.True(t, getAssociatedTagCount() > 0, "Database should have associated tags to delete")

		getEvidenceCount := makeDBRowCounter(t, db, "evidence", "uuid=?", i.EvidenceUUID)
		require.Equal(t, int64(1), getEvidenceCount(), "Database should have evidence to delete")

		ctx := contextForUser(UserRon, db)
		err := services.DeleteEvidence(ctx, db, memStore, i)
		require.NoError(t, err)
		require.Equal(t, int64(0), getEvidenceCount(), "Database should have deleted the evidence")
		require.Equal(t, int64(0), getAssociatedTagCount(), "Database should have deleted associated tags")
		_, err = memStore.Read(contentStoreKey)
		require.Error(t, err)
	})
}

func TestDeleteEvidenceWithPropogation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserRon, db)
		memStore := createPopulatedMemStore(seed)

		masterEvidence := EviDobby
		i := services.DeleteEvidenceInput{
			OperationSlug:            OpChamberOfSecrets.Slug,
			EvidenceUUID:             masterEvidence.UUID,
			DeleteAssociatedFindings: true,
		}
		getAssociatedTagCount := makeDBRowCounter(t, db, "tag_evidence_map", "evidence_id=?", masterEvidence.ID)
		require.True(t, getAssociatedTagCount() > 0, "Database should have associated tags to delete")

		getEvidenceCount := makeDBRowCounter(t, db, "evidence", "uuid=?", masterEvidence.UUID)
		require.Equal(t, int64(1), getEvidenceCount(), "Database should have evidence to delete")

		getMappedFindingCount := makeDBRowCounter(t, db, "evidence_finding_map", "evidence_id=?", masterEvidence.ID)
		require.True(t, getMappedFindingCount() > 0, "Database should have some mapped finding to delete")

		associatedFindingIDs := getAssociatedFindings(t, db, masterEvidence.ID)
		require.True(t, len(associatedFindingIDs) > 0, "Database should have some associated finding to delete")

		err := services.DeleteEvidence(ctx, db, memStore, i)
		require.NoError(t, err)
		require.Equal(t, int64(0), getEvidenceCount(), "Database should have deleted the evidence")
		require.Equal(t, int64(0), getAssociatedTagCount(), "Database should have deleted evidence-to-tags mappings")
		require.Equal(t, int64(0), getMappedFindingCount(), "Database should have deleted evidence-to-findings mappings")
		postDeleteFindingIDs := []int64{}
		db.Select(&postDeleteFindingIDs, sq.Select("id").From("findings").Where(sq.Eq{"id": associatedFindingIDs}))
		require.Equal(t, []int64{}, postDeleteFindingIDs, "Associated findings should be removed")
	})
}

func getAssociatedFindings(t *testing.T, db *database.Connection, evidenceID int64) []int64 {
	query := sq.Select("finding_id").From("evidence_finding_map").
		Where(sq.Eq{"evidence_id": evidenceID})

	var rtn []int64
	err := db.Select(&rtn, query)

	require.Nil(t, err)
	return rtn
}

func TestListEvidenceForFinding(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)
		cs, _ := contentstore.NewMemStore()

		masterOp := OpChamberOfSecrets
		masterFinding := FindingBook2Magic
		allEvidence := getFullEvidenceByFindingID(t, db, masterFinding.ID)

		require.NotEqual(t, 0, len(allEvidence), "Some evidence should be present for this finding")

		input := services.ListEvidenceForFindingInput{
			OperationSlug: masterOp.Slug,
			FindingUUID:   FindingBook2Magic.UUID,
		}

		foundEvidence, err := services.ListEvidenceForFinding(ctx, db, cs, input)
		require.NoError(t, err)
		require.Equal(t, len(foundEvidence), len(allEvidence))
		validateEvidenceSets(t, foundEvidence, allEvidence, validateEvidence)
	})
}

func TestListEvidenceForOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)
		cs, _ := contentstore.NewMemStore()

		masterOp := OpChamberOfSecrets
		allEvidence := getFullEvidenceByOperationID(t, db, masterOp.ID)

		require.NotEqual(t, len(allEvidence), 0, "Some evidence should be present")

		input := services.ListEvidenceForOperationInput{
			OperationSlug: masterOp.Slug,
			Filters:       helpers.TimelineFilters{},
		}

		foundEvidence, err := services.ListEvidenceForOperation(ctx, db, cs, input)
		require.NoError(t, err)
		require.Equal(t, len(foundEvidence), len(allEvidence))
		validateEvidenceSets(t, toRealEvidenceList(foundEvidence), allEvidence, validateEvidence)
	})
}

func TestReadEvidence(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)
		cs, _ := contentstore.NewMemStore()

		masterOp := OpChamberOfSecrets
		masterEvidence := EviFlyingCar

		neitherInput := services.ReadEvidenceInput{
			OperationSlug: masterOp.Slug,
			EvidenceUUID:  masterEvidence.UUID,
		}

		retrievedEvidence, err := services.ReadEvidence(ctx, db, cs, neitherInput)
		require.NoError(t, err)
		validateReadEvidenceOutput(t, masterEvidence, retrievedEvidence)

		fullImg := []byte("full_image")
		thumbImg := []byte("thumb_image")
		thumbKey, err := cs.Upload(bytes.NewReader(thumbImg))
		require.NoError(t, err)
		fullKey, err := cs.Upload(bytes.NewReader(fullImg))
		require.NoError(t, err)

		err = db.Update(sq.Update("evidence").
			SetMap(map[string]interface{}{
				"full_image_key":  fullKey,
				"thumb_image_key": thumbKey,
			}).
			Where(sq.Eq{"id": masterEvidence.ID}))
		require.NoError(t, err)
		masterEvidence.FullImageKey = fullKey
		masterEvidence.ThumbImageKey = thumbKey

		retrievedEvidence, err = services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{
			OperationSlug: masterOp.Slug,
			EvidenceUUID:  masterEvidence.UUID,
			LoadPreview:   true,
		})
		require.NoError(t, err)
		validateReadEvidenceOutput(t, masterEvidence, retrievedEvidence)
		previewBytes, err := io.ReadAll(retrievedEvidence.Preview)
		require.NoError(t, err)
		require.Equal(t, thumbImg, previewBytes)

		retrievedEvidence, err = services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{
			OperationSlug: masterOp.Slug,
			EvidenceUUID:  masterEvidence.UUID,
			LoadMedia:     true,
		})
		require.NoError(t, err)
		validateReadEvidenceOutput(t, masterEvidence, retrievedEvidence)
		mediaBytes, err := io.ReadAll(retrievedEvidence.Media)
		require.NoError(t, err)
		require.Equal(t, fullImg, mediaBytes)
	})
}

func TestUpdateEvidence(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserRon, db)
		cs, _ := contentstore.NewMemStore()

		// tests for common fields
		masterOp := OpChamberOfSecrets
		masterEvidence := EviFlyingCar
		initialTags := seed.TagsForEvidence(masterEvidence)
		tagToAdd := TagMercury
		tagToRemove := TagSaturn
		description := "New Description"
		input := services.UpdateEvidenceInput{
			OperationSlug: masterOp.Slug,
			EvidenceUUID:  masterEvidence.UUID,
			Description:   &description,
			TagsToRemove:  []int64{tagToRemove.ID},
			TagsToAdd:     []int64{tagToAdd.ID},
		}
		require.Contains(t, initialTags, tagToRemove)
		require.NotContains(t, initialTags, tagToAdd)

		err := services.UpdateEvidence(ctx, db, cs, input)
		require.NoError(t, err)
		evi, err := services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{OperationSlug: masterOp.Slug, EvidenceUUID: masterEvidence.UUID})
		require.NoError(t, err)
		require.Equal(t, *input.Description, evi.Description)
		expectedTagIDs := make([]int64, 0, len(initialTags))
		for _, t := range initialTags {
			if t != tagToRemove {
				expectedTagIDs = append(expectedTagIDs, t.ID)
			}
		}
		expectedTagIDs = append(expectedTagIDs, tagToAdd.ID)
		require.Equal(t, sorted(expectedTagIDs), sorted(getTagIDsFromEvidenceID(t, db, masterEvidence.ID)))

		// test for content

		codeblockEvidence := EviTomRiddlesDiary
		newContent := "stabbed_with_basilisk_fang = False\n\ndef is_alive():\n  return not stabbed_with_basilisk_fang\n"
		input = services.UpdateEvidenceInput{
			OperationSlug: masterOp.Slug,
			EvidenceUUID:  codeblockEvidence.UUID,
			Description:   &codeblockEvidence.Description, // Note: A quirk with UpdateEvidence is that it will always update the description, even if it is empty.
			Content:       bytes.NewReader([]byte(newContent)),
		}

		err = services.UpdateEvidence(ctx, db, cs, input)
		require.NoError(t, err)
		evi, err = services.ReadEvidence(ctx, db, cs, services.ReadEvidenceInput{
			OperationSlug: masterOp.Slug,
			EvidenceUUID:  codeblockEvidence.UUID,
			LoadMedia:     true,
			LoadPreview:   true,
		})
		require.NoError(t, err)
		mediaBytes, err := io.ReadAll(evi.Media)
		require.NoError(t, err)
		previewBytes, err := io.ReadAll(evi.Preview)
		require.NoError(t, err)
		require.Equal(t, mediaBytes, previewBytes, "Preview and Media content should be identical for codeblocks")
		require.Equal(t, []byte(newContent), previewBytes)

		updatedEvidence := getEvidenceByID(t, db, codeblockEvidence.ID)
		require.Equal(t, updatedEvidence.ThumbImageKey, updatedEvidence.FullImageKey)
		require.NotEqual(t, "", updatedEvidence.FullImageKey)
	})
}

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
	return helpers.Map(in, helpers.Ptr[dtos.Tag])
}

func toRealEvidenceList(in []*dtos.Evidence) []dtos.Evidence {
	return helpers.Map(in, func(v *dtos.Evidence) dtos.Evidence {
		return *v
	})
}

func validateReadEvidenceOutput(t *testing.T, expected models.Evidence, actual *services.ReadEvidenceOutput) {
	require.Equal(t, expected.UUID, actual.UUID)
	require.Equal(t, expected.Description, actual.Description)
	require.Equal(t, expected.ContentType, actual.ContentType)
	require.Equal(t, expected.OccurredAt.Round(time.Second).UTC(), actual.OccurredAt)
}

func validateInsertedEvidence(t *testing.T, evi *dtos.Evidence, src services.CreateEvidenceInput,
	user models.User, op models.Operation, store *contentstore.MemStore, db *database.Connection,
	rawContent []byte) {

	fullEvidence := getEvidenceByUUID(t, db, evi.UUID)
	require.Equal(t, src.Description, fullEvidence.Description)
	require.Equal(t, src.ContentType, fullEvidence.ContentType)
	require.Equal(t, user.ID, fullEvidence.OperatorID)
	require.Equal(t, op.ID, fullEvidence.OperationID, "Associated with right operation")
	if src.Content != nil {
		require.NotEqual(t, "", fullEvidence.FullImageKey, "Keys should be populated")
		require.NotEqual(t, "", fullEvidence.ThumbImageKey, "Keys should be populated")
		fullReader, _ := store.Read(fullEvidence.FullImageKey)
		fullContentBytes, _ := io.ReadAll(fullReader)
		require.Equal(t, rawContent, fullContentBytes)
	}

	tagIDs := getTagIDsFromEvidenceID(t, db, fullEvidence.ID)

	require.Equal(t, sorted(tagIDs), sorted(src.TagIDs))
}

func TestMoveEvidence(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		startingOp := OpChamberOfSecrets
		endingOp := OpSorcerersStone
		sourceEvidence := EviPetrifiedHermione //shares tags between the two operations

		input := services.MoveEvidenceInput{
			SourceOperationSlug: startingOp.Slug,
			TargetOperationSlug: endingOp.Slug,
			EvidenceUUID:        sourceEvidence.UUID,
		}

		// scenario 1: User present in both, cannot write dst [should fail]
		ctx := contextForUser(UserHermione, db)
		err := services.MoveEvidence(ctx, db, input)
		require.Error(t, err)

		// scenario 2: User present in both, cannot write src [should fail]
		ctx = contextForUser(UserSeamus, db)
		err = services.MoveEvidence(ctx, db, input)
		require.Error(t, err)

		// scenario 3: User present in src, cannot write dst [should fail]
		ctx = contextForUser(UserGinny, db)
		err = services.MoveEvidence(ctx, db, input)
		require.Error(t, err)

		// scenario 4: User present in dst, cannot write src [should fail]
		ctx = contextForUser(UserNeville, db)
		err = services.MoveEvidence(ctx, db, input)
		require.Error(t, err)

		// // scenario 5: User present in both, cannot write to both [should succeed]
		ctx = contextForUser(UserHarry, db)
		err = services.MoveEvidence(ctx, db, input)
		require.NoError(t, err)

		updatedEvidence := getEvidenceByUUID(t, db, sourceEvidence.UUID)
		require.Equal(t, updatedEvidence.OperationID, endingOp.ID)
		associatedTags := getTagIDsFromEvidenceID(t, db, updatedEvidence.ID)
		require.Equal(t, sorted(associatedTags), sorted([]int64{CommonTagWhoSS.ID, CommonTagWhatSS.ID}))
	})
}
