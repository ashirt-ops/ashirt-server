// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/services"
	"github.com/stretchr/testify/require"
)

func TestCreateOperation(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	ctx := fullContext(UserRon.ID, &policy.FullAccess{})

	i := services.CreateOperationInput{
		Slug:    "???",
		OwnerID: UserRon.ID,
		Name:    "Ron's Op",
	}
	createdOp, err := services.CreateOperation(ctx, db, i)
	require.Error(t, err)

	i = services.CreateOperationInput{
		Slug:    "rop",
		OwnerID: UserRon.ID,
		Name:    "Ron's Op",
	}
	createdOp, err = services.CreateOperation(ctx, db, i)
	require.NoError(t, err)
	fullOp := getOperationFromSlug(t, db, createdOp.Slug)

	require.NotEqual(t, 0, fullOp.ID)
	require.Equal(t, i.Name, fullOp.Name)
	require.Equal(t, models.OperationStatusPlanning, fullOp.Status, "status should default to 'Planning'")

	attachedUsers := getUserRolesForOperationByOperationID(t, db, fullOp.ID)
	require.Equal(t, 1, len(attachedUsers))
	require.Equal(t, policy.OperationRoleAdmin, attachedUsers[0].Role, "Creator of operation should have admin role for that operation")
	require.Equal(t, i.OwnerID, attachedUsers[0].UserID)
}

func TestSanitizeOperationSlug(t *testing.T) {
	require.Equal(t, services.SanitizeOperationSlug("?One?Two?Three?"), "one-two-three")
	require.Equal(t, services.SanitizeOperationSlug("Harry"), "harry")
	require.Equal(t, services.SanitizeOperationSlug("Harry Potter"), "harry-potter")
	require.Equal(t, services.SanitizeOperationSlug("fancy_name"), "fancy-name")
	require.Equal(t, services.SanitizeOperationSlug("Lots_Of-Fancy! Characters"), "lots-of-fancy-characters")
}
