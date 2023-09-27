// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/database/seeding"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"
)

func TestCreateGlobalVar(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserHarry, db)

		// verify non-admin cannot create var
		i := services.CreateGlobalVarInput{
			Name:  "Sectumsempra",
			Value: "slash a target",
		}
		_, err := services.CreateGlobalVar(ctx, db, i)
		require.Error(t, err)

		ctx = contextForUser(UserDumbledore, db)
		globalVar := VarExpelliarmus

		// verify name is invalid
		i = services.CreateGlobalVarInput{
			Name:  globalVar.Name,
			Value: "slash a target",
		}
		_, err = services.CreateGlobalVar(ctx, db, i)
		require.Error(t, err)

		// verify proper creation of a new var
		i = services.CreateGlobalVarInput{
			Name:  "Sectumsempra",
			Value: "slash a target",
		}
		createdGlobalVar, err := services.CreateGlobalVar(ctx, db, i)
		require.NoError(t, err)
		globalVar = getGlobalVarFromName(t, db, createdGlobalVar.Name)

		require.NotEqual(t, 0, globalVar.ID)
		require.Equal(t, i.Name, globalVar.Name)
		require.Equal(t, i.Value, globalVar.Value)
	})
}

func TestListGlobalVars(t *testing.T) {
	RunDisposableDBTestWithSeed(t, HarryPotterSeedData, func(db *database.Connection, _ TestSeedData) {
		// Verify that non-admins cannot list variables
		ctx := contextForUser(UserHarry, db)
		_, err := services.ListGlobalVars(ctx, db)
		require.Error(t, err)

		ctx = contextForUser(UserDumbledore, db)
		ops, err := services.ListGlobalVars(ctx, db)
		require.NoError(t, err)
		require.Equal(t, len(seeding.HarryPotterSeedData.GlobalVars), len(ops))
	})
}

func TestDeleteGlobalVar(t *testing.T) {
	RunDisposableDBTestWithSeed(t, HarryPotterSeedData, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserHarry, db)
		createdGlobalVar := VarExpelliarmus

		// Verify that non-admins cannot delete
		err := services.DeleteGlobalVar(ctx, db, createdGlobalVar.Name)
		require.Error(t, err)

		// Verify admins can delete
		ctx = contextForUser(UserDumbledore, db)
		err = services.DeleteGlobalVar(ctx, db, createdGlobalVar.Name)
		require.NoError(t, err)
	})
}

func TestUpdateGlobalVar(t *testing.T) {
	RunDisposableDBTestWithSeed(t, HarryPotterSeedData, func(db *database.Connection, _ TestSeedData) {
		initialVar := VarExpelliarmus

		// Verify that non-admins cannot update
		ctx := contextForUser(UserHarry, db)

		input := services.UpdateGlobalVarInput{
			GlobalVarName: initialVar.Name,
			NewName:       "Patronus",
			Value:         "Summon a Patronus",
		}

		err := services.UpdateGlobalVar(ctx, db, input)
		require.Error(t, err)

		// update name and value
		newVar := VarAscendio
		ctx = contextForUser(UserDumbledore, db)
		newName := "Accio"
		newValue := "Bring an object to you"

		input = services.UpdateGlobalVarInput{
			GlobalVarName: newVar.Name,
			NewName:       newName,
			Value:         newValue,
		}

		err = services.UpdateGlobalVar(ctx, db, input)
		require.NoError(t, err)

		updatedGlobalVar, err := services.LookupGlobalVar(db, newName)
		require.NoError(t, err)
		require.Equal(t, newName, updatedGlobalVar.Name)
		require.Equal(t, newValue, updatedGlobalVar.Value)

		// update only name
		newVar = VarImperio
		newName = "Expecto Patronum"
		input = services.UpdateGlobalVarInput{
			GlobalVarName: newVar.Name,
			NewName:       newName,
			Value:         "",
		}

		err = services.UpdateGlobalVar(ctx, db, input)
		require.NoError(t, err)
		updatedGlobalVar, err = services.LookupGlobalVar(db, newName)
		require.Equal(t, newName, updatedGlobalVar.Name)
		require.Equal(t, newVar.Value, updatedGlobalVar.Value)

		// update only value
		newVar = VarLumos
		newValue = "Summon a Patronus"
		input = services.UpdateGlobalVarInput{
			GlobalVarName: newVar.Name,
			NewName:       "",
			Value:         newValue,
		}

		err = services.UpdateGlobalVar(ctx, db, input)
		require.NoError(t, err)
		updatedGlobalVar, err = services.LookupGlobalVar(db, newVar.Name)
		require.Equal(t, newVar.Name, updatedGlobalVar.Name)
		require.Equal(t, newValue, updatedGlobalVar.Value)

		// Update name to another var that already exists
		newVar = VarAlohomora
		input = services.UpdateGlobalVarInput{
			GlobalVarName: newName,
			NewName:       newVar.Name,
			Value:         "",
		}

		err = services.UpdateGlobalVar(ctx, db, input)
		require.Error(t, err)
	})
}
