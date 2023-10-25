// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateOperationVarInput struct {
	OperationSlug string
	Name          string
	VarSlug       string
	Value         string
}

type UpdateOperationVarInput struct {
	Name          string
	Value         string
	VarSlug       string
	OperationSlug string
}

type DeleteOperationVarInput struct {
	Name string
}

func CreateOperationVar(ctx context.Context, db *database.Connection, i CreateOperationVarInput) (*dtos.OperationVar, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err := policy.Require(middleware.Policy(ctx), policy.CanCreateOpVars{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unable to create operation variable", backend.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}
	formattedName := helpers.StrToUpperCaseUnderscore(i.Name)

	cleanSlug := SanitizeSlug(i.VarSlug)
	if cleanSlug == "" {
		return nil, backend.BadInputErr(errors.New("Unable to create operation variable. Invalid operation variable slug"), "Slug must contain english letters or numbers")
	}
	formattedSlug := helpers.StrToLowerCaseUnderscore(cleanSlug)

	var varID int64

	listOfVarsInOperation, err := ListOperationVars(ctx, db, i.OperationSlug)
	for _, varInOperation := range listOfVarsInOperation {
		if varInOperation.Name == formattedName {
			return nil, backend.BadInputErr(errors.New("Unable to create operation variable. Invalid operation variable name"), "A variable with this name already exists in the operation")
		}
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		varID, _ = tx.Insert("operation_vars", map[string]interface{}{
			"name":  formattedName,
			"value": i.Value,
			"slug":  formattedSlug,
		})
		tx.Insert("var_operation_map", map[string]interface{}{
			"var_id":       varID,
			"operation_id": operation.ID,
		})
	})
	if err != nil {
		var errMessage string
		if database.IsAlreadyExistsError(err) {
			errMessage = "An operation variable with this name already exists"
		} else if database.InputIsTooLongError(err) {
			errMessage = "The variable name must be 255 characters or less"
		} else {
			errMessage = "Unable to add new operation variable"
		}

		return nil, backend.BadInputErr(backend.WrapError(errMessage, backend.DatabaseErr(err)), errMessage)
	}

	return &dtos.OperationVar{
		Name:    formattedName,
		Value:   i.Value,
		VarSlug: formattedSlug,
	}, nil
}

func DeleteOperationVar(ctx context.Context, db *database.Connection, varSlug string, operationSlug string) error {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return backend.WrapError("Unable to read operation", backend.UnauthorizedReadErr(err))
	}
	if err := policyRequireWithAdminBypass(ctx, policy.CanDeleteOpVars{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete operation variable", backend.UnauthorizedWriteErr(err))
	}

	err = db.Delete(sq.Delete("operation_vars").Where(sq.Eq{"slug": varSlug}))
	if err != nil {
		return backend.WrapError("Cannot delete operation variable", backend.DatabaseErr(err))
	}

	return nil
}

func ListOperationVars(ctx context.Context, db *database.Connection, operationSlug string) ([]*dtos.OperationVar, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to read operation", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanViewOpVars{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list operation variables", backend.UnauthorizedReadErr(err))
	}

	var operationVars = make([]models.OperationVar, 0)

	err = db.Select(&operationVars, sq.
		Select("ov.*").
		From("operations o").
		Join("var_operation_map vom ON o.id = vom.operation_id").
		Join("operation_vars ov ON ov.id = vom.var_id").
		Where(sq.Eq{"o.id": operation.ID}).
		OrderBy("ov.name ASC"))

	if err != nil {
		return nil, backend.WrapError("Cannot list operation variables", backend.DatabaseErr(err))
	}

	var operationVarsDTO = make([]*dtos.OperationVar, len(operationVars))
	for i, operationVar := range operationVars {
		operationVarsDTO[i] = &dtos.OperationVar{
			Name:    operationVar.Name,
			Value:   operationVar.Value,
			VarSlug: operationVar.Slug,
		}
	}

	return operationVarsDTO, nil
}

func UpdateOperationVar(ctx context.Context, db *database.Connection, i UpdateOperationVarInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to read operation", backend.UnauthorizedReadErr(err))
	}
	operationVar, err := LookupOperationVar(db, i.VarSlug)
	if err != nil {
		return backend.WrapError("Unable to update operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyOpVars{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to update operation", backend.UnauthorizedWriteErr(err))
	}
	formattedName := helpers.StrToUpperCaseUnderscore(i.Name)

	listOfVarsInOperation, err := ListOperationVars(ctx, db, i.OperationSlug)
	for _, varInOperation := range listOfVarsInOperation {
		if varInOperation.Name == formattedName {
			return backend.BadInputErr(errors.New("Unable to update operation variable. Invalid operation variable name"), "A variable with this name already exists in the operation")
		}
	}

	var val string
	var name string

	if i.Value != "" {
		val = i.Value
	} else {
		val = operationVar.Value
	}

	if i.Name != "" {
		name = formattedName
	} else {
		name = operationVar.Name
	}

	err = db.Update(sq.Update("operation_vars").
		SetMap(map[string]interface{}{
			"value": val,
			"name":  name,
		}).
		Where(sq.Eq{"id": operationVar.ID}))
	if err != nil {
		var errMessage string
		if database.InputIsTooLongErrorSq(err) {
			errMessage = "The variable name must be 255 characters or less"
		} else {
			errMessage = "Unable to update new operation variable"
		}

		return backend.BadInputErr(backend.WrapError(errMessage, backend.DatabaseErr(err)), errMessage)
	}

	return nil
}
