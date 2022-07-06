// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/services"

	sq "github.com/Masterminds/squirrel"
)

func TestCreateOperation(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

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
	require.Equal(t, models.OperationStatusPlanning, fullOp.Status, "status should default to 'Planning'")

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
}

func TestDeleteOperation(t *testing.T) {
	db := initTest(t)
	defer db.DB.Close()
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserHarry.ID, &policy.Deny{})
	memStore := createPopulatedMemStore(HarryPotterSeedData)

	masterOp := OpChamberOfSecrets
	originalEvidence := getEvidenceForOperation(t, db, masterOp.ID)

	// Verify that non-admins cannot delete
	err := services.DeleteOperation(ctx, db, memStore, masterOp.Slug)
	require.Error(t, err)

	// Verify admins can delete
	ctx = fullContext(UserRon.ID, &policy.FullAccess{})
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
}

func TestListOperations(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	validateOperationList := func(receivedOps []*dtos.Operation, expectedOps []models.Operation) {
		for _, op := range receivedOps {
			var expected *models.Operation = nil
			for _, fOp := range expectedOps {
				if fOp.Slug == op.Slug {
					expected = &fOp
					break
				}
			}
			require.NotNil(t, expected, "Result should have matching value")
			validateOp(t, *expected, op)
		}
	}

	normalUser := UserRon
	expectedOps := getOperationsForUser(t, db, normalUser.ID)

	ops, err := services.ListOperations(contextForUser(normalUser, db), db)
	require.NoError(t, err)
	require.Equal(t, len(ops), len(expectedOps))
	validateOperationList(ops, expectedOps)

	// validate headless users
	headlessUser := UserHeadlessNick
	fullOps := getOperations(t, db)

	ops, err = services.ListOperations(contextForUser(headlessUser, db), db)
	require.NoError(t, err)
	require.Equal(t, len(ops), len(fullOps))
	validateOperationList(ops, fullOps)
}

func TestListOperationsForAdmin(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContextAsAdmin(UserDumbledore.ID, &policy.FullAccess{})

	fullOps := getOperations(t, db)
	require.NotEqual(t, len(fullOps), 0, "Some number of operations should exist")

	ops, err := services.ListOperationsForAdmin(ctx, db)
	require.NoError(t, err)
	require.Equal(t, len(ops), len(fullOps))
	for _, op := range ops {
		var expected *models.Operation = nil
		for _, fOp := range fullOps {
			if fOp.Slug == op.Slug {
				expected = &fOp
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validateOp(t, *expected, op)
	}

	// verify non admins don't have access

	ctx = fullContext(UserDraco.ID, &policy.FullAccess{}) // Note: not an admin
	ops, err = services.ListOperationsForAdmin(ctx, db)
	require.Error(t, err)
	require.Equal(t, "Requesting user is not an admin", err.Error())
}

func TestSanitizeOperationSlug(t *testing.T) {
	require.Equal(t, services.SanitizeOperationSlug("?One?Two?Three?"), "one-two-three")
	require.Equal(t, services.SanitizeOperationSlug("Harry"), "harry")
	require.Equal(t, services.SanitizeOperationSlug("Harry Potter"), "harry-potter")
	require.Equal(t, services.SanitizeOperationSlug("fancy_name"), "fancy-name")
	require.Equal(t, services.SanitizeOperationSlug("Lots_Of-Fancy! Characters"), "lots-of-fancy-characters")
}

func TestUpdateOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	// tests for common fields
	masterOp := OpChamberOfSecrets
	input := services.UpdateOperationInput{
		OperationSlug: masterOp.Slug,
		Name:          "New Name",
		Status:        models.OperationStatusComplete,
	}
	require.NotEqual(t, masterOp.Status, input.Status)

	err := services.UpdateOperation(ctx, db, input)
	require.NoError(t, err)
	updatedOperation, err := services.ReadOperation(ctx, db, masterOp.Slug)
	require.NoError(t, err)
	require.Equal(t, input.Name, updatedOperation.Name)
	require.Equal(t, input.Status, updatedOperation.Status)
}


func TestReadOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	masterOp := OpChamberOfSecrets

	retrievedOp, err := services.ReadOperation(ctx, db, masterOp.Slug)
	require.NoError(t, err)

	require.Equal(t, masterOp.Slug, retrievedOp.Slug)
	require.Equal(t, masterOp.Name, retrievedOp.Name)
	require.Equal(t, masterOp.Status, retrievedOp.Status)
	require.Equal(t, len(HarryPotterSeedData.UsersForOp(masterOp)), retrievedOp.NumUsers)
}


func validateOp(t *testing.T, expected models.Operation, actual *dtos.Operation) {
	require.Equal(t, expected.Slug, actual.Slug, "Slugs should match")
	require.Equal(t, expected.Name, actual.Name, "Names should match")
	require.Equal(t, expected.Status, actual.Status, "Status should match")
}
