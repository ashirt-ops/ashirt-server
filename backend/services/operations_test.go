package services_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

func TestCreateOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		// verify slug name is invalid
		i := services.CreateOperationInput{
			Slug:    "???",
			OwnerID: UserRon.ID,
			Name:    "Ron's Op",
		}
		_, err := services.CreateOperation(ctx, db, i)
		require.Error(t, err)

		// verify proper creation of a new operation
		i = services.CreateOperationInput{
			Slug:    "rop",
			OwnerID: UserRon.ID,
			Name:    "Ron's Op",
		}
		createdOp, err := services.CreateOperation(ctx, db, i)
		require.NoError(t, err)
		fullOp := getOperationFromSlug(t, db, createdOp.Slug)

		require.NotEqual(t, 0, fullOp.ID)
		require.Equal(t, i.Name, fullOp.Name)

		attachedUsers := getUserRolesForOperationByOperationID(t, db, fullOp.ID)
		require.Equal(t, 1, len(attachedUsers))
		require.Equal(t, policy.OperationRoleAdmin, attachedUsers[0].Role, "Creator of operation should have admin role for that operation")
		require.Equal(t, i.OwnerID, attachedUsers[0].UserID)

		attachedTags := getTagFromOperationID(t, db, fullOp.ID)
		defaultTags := getDefaultTags(t, db)
		expectedTags := make([]models.Tag, len(defaultTags))
		for idx, tag := range defaultTags {
			expectedTags[idx].ColorName = tag.ColorName
			expectedTags[idx].Name = tag.Name
		}

		for _, tag := range attachedTags {
			foundIndex := -1
			for idx, eTag := range expectedTags {
				if tag.Name == eTag.Name && tag.ColorName == eTag.ColorName {
					foundIndex = idx
				}
			}
			require.NotEqual(t, -1, foundIndex, "Each of the created tags must be from default tags")
			expectedTags = append(expectedTags[:foundIndex], expectedTags[foundIndex+1:]...)
		}
		require.Empty(t, expectedTags, "All of the expected tags must be used")
	})
}

func TestDeleteOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserHarry, db)
		memStore := createPopulatedMemStore(seed)

		masterOp := OpChamberOfSecrets
		originalEvidence := getEvidenceForOperation(t, db, masterOp.ID)

		// Verify that non-admins cannot delete
		err := services.DeleteOperation(ctx, db, memStore, masterOp.Slug)
		require.Error(t, err)

		// Verify admins can delete
		ctx = contextForUser(UserRon, db)
		err = services.DeleteOperation(ctx, db, memStore, masterOp.Slug)
		require.NoError(t, err)
		// ensure content was removed
		for _, evi := range originalEvidence {
			_, err = memStore.Read(evi.FullImageKey)
			require.Error(t, err)
			_, err = memStore.Read(evi.ThumbImageKey)
			require.Error(t, err)
		}
		var dbOp models.Operation
		err = db.Get(&dbOp, sq.Select("*").From("operations").Where(sq.Eq{"id": masterOp.ID}))
		// assuming that if this row was deleted, then all other rows must have been deleted (via foreign key constraint)
		require.Error(t, err)

		// Verify Super admins can delete
		// TODO
	})
}

func TestListOperations(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		validateOperationList := func(receivedOps []*dtos.Operation, expectedOps []*dtos.Operation) {
			for _, op := range receivedOps {
				var expected *dtos.Operation = nil
				for _, fOp := range expectedOps {
					if fOp.Slug == op.Slug {
						expected = fOp
						break
					}
				}
				require.NotNil(t, expected, "Result should have matching value")
				validateOp(t, expected, op)
			}
		}

		normalUser := UserRon
		expectedOps := getOperationsForUser(t, db, normalUser)

		ops, err := services.ListOperations(contextForUser(normalUser, db), db)
		require.NoError(t, err)
		require.Equal(t, len(expectedOps), len(ops))
		validateOperationList(ops, expectedOps)

		opsAndPrefs := getFavoritesByUserID(t, db, normalUser.ID)
		for _, opPrefs := range opsAndPrefs {
			_, found := helpers.Find(ops, func(expectedOp *dtos.Operation) bool {
				return (*expectedOp).Slug == opPrefs.Slug
			})
			require.NotNil(t, found)
			fav := (**found).Favorite
			require.Equal(t, fav, opPrefs.IsFavorite)
		}

		// validate headless users
		headlessUser := UserHeadlessNick
		fullOps := getOperationsForUser(t, db, headlessUser)

		ops, err = services.ListOperations(contextForUser(headlessUser, db), db)
		require.NoError(t, err)
		require.Equal(t, len(ops), len(fullOps))
		validateOperationList(ops, fullOps)
	})
}

func TestListOperationsForAdmin(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserDumbledore, db)

		fullOps := getOperations(t, db)
		require.NotEqual(t, len(fullOps), 0, "Some number of operations should exist")

		ops, err := services.ListOperationsForAdmin(ctx, db)
		require.NoError(t, err)
		require.Equal(t, len(ops), len(fullOps))
		for _, op := range ops {
			var expected *dtos.Operation = nil
			for _, fOp := range ops {
				if fOp.Slug == op.Slug {
					expected = fOp
					break
				}
			}
			require.NotNil(t, expected, "Result should have matching value")
			validateOp(t, expected, op)
		}

		// verify non admins don't have access
		ctx = contextForUser(UserDraco, db)
		_, err = services.ListOperationsForAdmin(ctx, db)
		require.Error(t, err)
		require.Equal(t, "Requesting user is not an admin", err.Error())
	})
}

func TestSetFavoriteOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		slug := OpGobletOfFire.Slug

		isFavorite := getFavoriteForOperation(t, db, slug, normalUser.ID)
		require.Equal(t, isFavorite, false)

		i := services.SetFavoriteInput{slug, true}
		err := services.SetFavoriteOperation(contextForUser(normalUser, db), db, i)
		require.NoError(t, err)

		isFavorite = getFavoriteForOperation(t, db, slug, normalUser.ID)
		require.Equal(t, isFavorite, true)
	})
}

func TestUpdateOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		// tests for common fields
		masterOp := OpChamberOfSecrets
		input := services.UpdateOperationInput{
			OperationSlug: masterOp.Slug,
			Name:          "New Name",
		}

		err := services.UpdateOperation(ctx, db, input)
		require.NoError(t, err)
		updatedOperation, err := services.ReadOperation(ctx, db, masterOp.Slug)
		require.NoError(t, err)
		require.Equal(t, input.Name, updatedOperation.Name)
	})
}

func TestReadOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpChamberOfSecrets
		originalEvidence := getEvidenceForOperation(t, db, masterOp.ID)
		harEvidence := helpers.Filter(originalEvidence, func(evi models.Evidence) bool { return evi.ContentType == "http-request-cycle" })

		retrievedOp, err := services.ReadOperation(ctx, db, masterOp.Slug)
		require.NoError(t, err)

		require.Equal(t, masterOp.Slug, retrievedOp.Slug)
		require.Equal(t, masterOp.Name, retrievedOp.Name)
		require.Equal(t, 6, retrievedOp.NumUsers)
		require.Equal(t, true, retrievedOp.Favorite)
		require.Equal(t, len(originalEvidence), retrievedOp.NumEvidence)
		require.Equal(t, 12, retrievedOp.NumTags)
		require.Equal(t, 1, len(retrievedOp.TopContribs))
		require.Equal(t, "harry.potter", retrievedOp.TopContribs[0].Slug)
		require.Equal(t, int64(2), retrievedOp.EvidenceCount.CodeblockCount)
		require.Equal(t, int64(6), retrievedOp.EvidenceCount.ImageCount)
		require.Equal(t, int64(0), retrievedOp.EvidenceCount.RecordingCount)
		require.Equal(t, int64(0), retrievedOp.EvidenceCount.EventCount)
		require.Equal(t, int64(len(harEvidence)), retrievedOp.EvidenceCount.HarCount)

		require.Equal(t, len(seed.UsersForOp(masterOp)), retrievedOp.NumUsers)
	})
}

func validateOp(t *testing.T, expected *dtos.Operation, actual *dtos.Operation) {
	require.Equal(t, expected.Slug, actual.Slug, "Slugs should match")
	require.Equal(t, expected.Name, actual.Name, "Names should match")
	require.Equal(t, expected.Favorite, actual.Favorite, "Favorite should match")
	require.Equal(t, expected.NumUsers, actual.NumUsers, "NumUsers should match")
	require.Equal(t, expected.NumEvidence, actual.NumEvidence, "NumEvidence should match")
	require.Equal(t, expected.EvidenceCount, actual.EvidenceCount, "EvidenceCount should match")
	require.Equal(t, expected.TopContribs, actual.TopContribs, "TopContribs should match")
	require.Equal(t, expected.NumTags, actual.NumTags, "NumTags should match")
}
