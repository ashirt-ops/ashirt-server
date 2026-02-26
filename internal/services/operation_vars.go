package services

import (
	"context"
	stderrors "errors"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/dtos"
	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/server/middleware"

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
		return nil, errors.WrapError("Unable to create operation variable", errors.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, errors.MissingValueErr("Name")
	}
	formattedName := helpers.StrToUpperCaseUnderscore(i.Name)

	cleanSlug := SanitizeSlug(i.VarSlug)
	if cleanSlug == "" {
		return nil, errors.BadInputErr(stderrors.New("Unable to create operation variable. Invalid operation variable slug"), "Slug must contain english letters or numbers")
	}
	formattedSlug := helpers.StrToLowerCaseUnderscore(cleanSlug)

	var varID int64

	listOfVarsInOperation, err := ListOperationVars(ctx, db, i.OperationSlug)
	for _, varInOperation := range listOfVarsInOperation {
		if varInOperation.Name == formattedName {
			return nil, errors.BadInputErr(stderrors.New("Unable to create operation variable. Invalid operation variable name"), "A variable with this name already exists in the operation")
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

		return nil, errors.BadInputErr(errors.WrapError(errMessage, errors.DatabaseErr(err)), errMessage)
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
		return errors.WrapError("Unable to read operation", errors.UnauthorizedReadErr(err))
	}
	if err := policyRequireWithAdminBypass(ctx, policy.CanDeleteOpVars{OperationID: operation.ID}); err != nil {
		return errors.WrapError("Unwilling to delete operation variable", errors.UnauthorizedWriteErr(err))
	}

	err = db.Delete(sq.Delete("operation_vars").Where(sq.Eq{"slug": varSlug}))
	if err != nil {
		return errors.WrapError("Cannot delete operation variable", errors.DatabaseErr(err))
	}

	return nil
}

func ListOperationVars(ctx context.Context, db *database.Connection, operationSlug string) ([]*dtos.OperationVar, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, errors.WrapError("Unable to read operation", errors.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanViewOpVars{OperationID: operation.ID}); err != nil {
		return nil, errors.WrapError("Unwilling to list operation variables", errors.UnauthorizedReadErr(err))
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
		return nil, errors.WrapError("Cannot list operation variables", errors.DatabaseErr(err))
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
		return errors.WrapError("Unable to read operation", errors.UnauthorizedReadErr(err))
	}
	operationVar, err := LookupOperationVar(db, i.VarSlug)
	if err != nil {
		return errors.WrapError("Unable to update operation", errors.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyOpVars{OperationID: operation.ID}); err != nil {
		return errors.WrapError("Unwilling to update operation", errors.UnauthorizedWriteErr(err))
	}
	formattedName := helpers.StrToUpperCaseUnderscore(i.Name)

	listOfVarsInOperation, err := ListOperationVars(ctx, db, i.OperationSlug)
	for _, varInOperation := range listOfVarsInOperation {
		if varInOperation.Name == formattedName {
			return errors.BadInputErr(stderrors.New("Unable to update operation variable. Invalid operation variable name"), "A variable with this name already exists in the operation")
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

		return errors.BadInputErr(errors.WrapError(errMessage, errors.DatabaseErr(err)), errMessage)
	}

	return nil
}
