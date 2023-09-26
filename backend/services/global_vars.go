// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateGlobalVarInput struct {
	Name    string
	OwnerID int64
	Value   string
}

type UpdateGlobalVarInput struct {
	GlobalVarName string
	Value         string
	NewName       string
}

type DeleteGlobalVarInput struct {
	Name string
}

func CreateGlobalVar(ctx context.Context, db *database.Connection, i CreateGlobalVarInput) (*dtos.GlobalVar, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.CanCreateGlobalVars{}); err != nil {
		return nil, backend.WrapError("Unable to create global variable", backend.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}

	// TODO TN can this be empty? I think it should be able to be
	// if i.Value == "" {
	// 	return nil, backend.MissingValueErr("Value")
	// }

	// todo TN does this need to happen?
	// cleanSlug := SanitizeSlug(i.Name)
	// if cleanSlug == "" {
	// 	return nil, backend.BadInputErr(errors.New("Unable to create operation. Invalid operation slug"), "Slug must contain english letters or numbers")
	// }

	globalVarID, err := db.Insert("global_vars", map[string]interface{}{
		"name":  i.Name,
		"value": i.Value,
	})
	if err != nil {
		fmt.Println("other error", err)
		if database.IsAlreadyExistsError(err) {
			return nil, backend.BadInputErr(backend.WrapError("global variable already exists", err), "A global variable with this name already exists")
		}
		return nil, backend.WrapError("Unable to add new global variable", backend.DatabaseErr(err))
	}

	return &dtos.GlobalVar{
		ID:    globalVarID,
		Name:  i.Name,
		Value: i.Value,
	}, nil
}

func DeleteGlobalVar(ctx context.Context, db *database.Connection, name string) error {
	globalVar, err := LookupGlobalVar(db, name)
	if err != nil {
		return backend.WrapError("Unable to delete global variable", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanDeleteGlobalVar{GlobalVarID: globalVar.ID}); err != nil {
		return backend.WrapError("Unwilling to delete global variable", backend.UnauthorizedWriteErr(err))
	}

	err = db.Delete(sq.Delete("global_vars").Where(sq.Eq{"name": name}))
	if err != nil {
		return backend.WrapError("Cannot delete global variable", backend.DatabaseErr(err))
	}

	return nil
}

// ListQueriesForOperation retrieves all saved queries for a given globalVar id
func ListGlobalVars(ctx context.Context, db *database.Connection) ([]*dtos.GlobalVar, error) {

	// TODO TN - what needs to happen here to check this?
	// if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: globalVar.ID}); err != nil {
	// 	return nil, backend.WrapError("Unwilling to list global variables", backend.UnauthorizedReadErr(err))
	// }

	var globalVars = make([]models.GlobalVar, 0)
	err := db.Select(&globalVars, sq.Select("*").
		From("global_vars").
		// TODO TN - do we want to do this?
		OrderBy("name ASC"))

	if err != nil {
		return nil, backend.WrapError("Cannot list global variables", backend.DatabaseErr(err))
	}

	var globalVarsDTO = make([]*dtos.GlobalVar, len(globalVars))
	for i, globalVar := range globalVars {
		globalVarsDTO[i] = &dtos.GlobalVar{
			ID:    globalVar.ID,
			Name:  globalVar.Name,
			Value: globalVar.Value,
		}
	}

	return globalVarsDTO, nil
}

func UpdateGlobalVar(ctx context.Context, db *database.Connection, i UpdateGlobalVarInput) error {
	globalVar, err := LookupGlobalVar(db, i.GlobalVarName)
	if err != nil {
		return backend.WrapError("Unable to update operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyGlobalVar{GlobalVarID: globalVar.ID}); err != nil {
		return backend.WrapError("Unwilling to update operation", backend.UnauthorizedWriteErr(err))
	}

	var val string
	var name string
	// TODO TN test this to make sure it works as intended
	if i.Value != "" {
		val = i.Value
	} else {
		val = globalVar.Value
	}

	if i.NewName != "" {
		name = i.NewName
	} else {
		name = globalVar.Name
	}

	err = db.Update(sq.Update("global_vars").
		SetMap(map[string]interface{}{
			"value": val,
			"name":  name,
		}).
		Where(sq.Eq{"id": globalVar.ID}))
	if err != nil {
		if database.IsAlreadyExistsErrorSq(err) {
			return backend.BadInputErr(backend.WrapError("Global variable already exists", err), "A global variable with this name already exists")
		}
		return backend.WrapError("Cannot update global variable", backend.DatabaseErr(err))
	}

	return nil
}
