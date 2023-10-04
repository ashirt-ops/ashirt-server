// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"
)

func TestCreateOperationVar(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		// verify non-admin cannot create var
		i := services.CreateOperationVarInput{
			OperationSlug: OpSorcerersStone.Slug,
			VarSlug:       "Sectumsempra",
			Name:          "Sectumsempra",
			Value:         "slash a target",
		}
		_, err := services.CreateOperationVar(ctx, db, i)
		require.Error(t, err)

		ctx = contextForUser(UserHarry, db)
		operationVar := OpVarImmobulus

		// verify name is invalid
		i = services.CreateOperationVarInput{
			OperationSlug: OpSorcerersStone.Slug,
			VarSlug:       operationVar.Slug,
			Name:          "Sectumsempra",
			Value:         "slash a target",
		}
		_, err = services.CreateOperationVar(ctx, db, i)
		require.Error(t, err)

		// verify proper creation of a new var
		i = services.CreateOperationVarInput{
			OperationSlug: OpSorcerersStone.Slug,
			VarSlug:       "Sectumsempra",
			Name:          "Sectumsempra",
			Value:         "slash a target",
		}
		createdOperationVar, err := services.CreateOperationVar(ctx, db, i)
		require.NoError(t, err)
		operationVar = getOperationVarFromSlug(t, db, createdOperationVar.VarSlug)

		require.NotEqual(t, 0, operationVar.ID)
		require.Equal(t, i.VarSlug, operationVar.Slug)
		require.Equal(t, i.Name, operationVar.Name)
		require.Equal(t, i.Value, operationVar.Value)
	})
}

func TestListOperationVars(t *testing.T) {
	RunDisposableDBTestWithSeed(t, HarryPotterSeedData, func(db *database.Connection, _ TestSeedData) {
		// Verify that non-admins cannot list variables
		ctx := contextForUser(UserRon, db)
		_, err := services.ListOperationVars(ctx, db, OpSorcerersStone.Slug)
		require.Error(t, err)

		ctx = contextForUser(UserHarry, db)
		opVars, err := services.ListOperationVars(ctx, db, OpSorcerersStone.Slug)
		require.NoError(t, err)
		require.Equal(t, 2, len(opVars))
	})
}

func TestDeleteOperationVar(t *testing.T) {
	RunDisposableDBTestWithSeed(t, HarryPotterSeedData, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserRon, db)
		operationVar := OpVarObscuro

		// Verify that non-admins cannot delete
		err := services.DeleteOperationVar(ctx, db, operationVar.Slug, OpSorcerersStone.Slug)
		require.Error(t, err)

		// Verify admins can delete
		ctx = contextForUser(UserHarry, db)
		err = services.DeleteOperationVar(ctx, db, operationVar.Slug, OpSorcerersStone.Slug)
		require.NoError(t, err)
	})
}

func TestUpdateOperationVar(t *testing.T) {
	RunDisposableDBTestWithSeed(t, HarryPotterSeedData, func(db *database.Connection, _ TestSeedData) {
		initialVar := OpVarProtego

		// Verify that non-admins cannot update
		ctx := contextForUser(UserHarry, db)

		input := services.UpdateOperationVarInput{
			VarSlug:       initialVar.Slug,
			Name:          "Patronus",
			Value:         "Summon a Patronus",
			OperationSlug: OpChamberOfSecrets.Slug,
		}

		err := services.UpdateOperationVar(ctx, db, input)
		require.Error(t, err)

		// update name and value
		newVar := OpVarReparo
		ctx = contextForUser(UserRon, db)
		newName := "Accio"
		newValue := "Bring an object to you"

		input = services.UpdateOperationVarInput{
			VarSlug:       newVar.Slug,
			OperationSlug: OpChamberOfSecrets.Slug,
			Name:          newName,
			Value:         newValue,
		}

		err = services.UpdateOperationVar(ctx, db, input)
		require.NoError(t, err)

		updatedOperationVar, err := services.LookupOperationVar(db, newVar.Slug)
		require.NoError(t, err)
		require.Equal(t, newName, updatedOperationVar.Name)
		require.Equal(t, newValue, updatedOperationVar.Value)
		require.Equal(t, newVar.Slug, updatedOperationVar.Slug)

		ctx = contextForUser(UserHarry, db)
		// update only name
		newVar = OpVarStupefy
		newName = "Expecto Patronum"
		input = services.UpdateOperationVarInput{
			VarSlug:       newVar.Slug,
			Name:          newName,
			Value:         "",
			OperationSlug: OpGobletOfFire.Slug,
		}

		err = services.UpdateOperationVar(ctx, db, input)
		require.NoError(t, err)
		updatedOperationVar, err = services.LookupOperationVar(db, newVar.Slug)
		require.Equal(t, newName, updatedOperationVar.Name)
		require.Equal(t, newVar.Value, updatedOperationVar.Value)
		require.Equal(t, newVar.Slug, updatedOperationVar.Slug)

		// update only value
		newVar = OpVarWingardiumLeviosa
		newValue = "Summon a Patronus"
		input = services.UpdateOperationVarInput{
			VarSlug:       newVar.Slug,
			Name:          "",
			Value:         newValue,
			OperationSlug: OpGobletOfFire.Slug,
		}

		err = services.UpdateOperationVar(ctx, db, input)
		require.NoError(t, err)
		updatedOperationVar, err = services.LookupOperationVar(db, newVar.Slug)
		require.Equal(t, newVar.Name, updatedOperationVar.Name)
		require.Equal(t, newValue, updatedOperationVar.Value)
	})
}
