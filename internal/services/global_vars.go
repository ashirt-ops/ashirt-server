package services

import (
	"context"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/dtos"
	"github.com/ashirt-ops/ashirt-server/internal/errorwrap"
	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateGlobalVarInput struct {
	Name    string
	OwnerID int64
	Value   string
}

type UpdateGlobalVarInput struct {
	Name    string
	Value   string
	NewName string
}

type DeleteGlobalVarInput struct {
	Name string
}

func CreateGlobalVar(ctx context.Context, db *database.Connection, i CreateGlobalVarInput) (*dtos.GlobalVar, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, errorwrap.WrapError("Unable to create global variable", errorwrap.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, errorwrap.MissingValueErr("Name")
	}
	formattedName := helpers.StrToUpperCaseUnderscore(i.Name)

	_, err := db.Insert("global_vars", map[string]interface{}{
		"name":  formattedName,
		"value": i.Value,
	})
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return nil, errorwrap.BadInputErr(errorwrap.WrapError("global variable already exists", err), "A global variable with this name already exists")
		}
		return nil, errorwrap.WrapError("Unable to add new global variable", errorwrap.DatabaseErr(err))
	}

	return &dtos.GlobalVar{
		Name:  formattedName,
		Value: i.Value,
	}, nil
}

func DeleteGlobalVar(ctx context.Context, db *database.Connection, name string) error {
	if err := policyRequireWithAdminBypass(ctx, policy.AdminUsersOnly{}); err != nil {
		return errorwrap.WrapError("Unwilling to delete global variable", errorwrap.UnauthorizedWriteErr(err))
	}

	err := db.Delete(sq.Delete("global_vars").Where(sq.Eq{"name": name}))
	if err != nil {
		return errorwrap.WrapError("Cannot delete global variable", errorwrap.DatabaseErr(err))
	}

	return nil
}

func ListGlobalVars(ctx context.Context, db *database.Connection) ([]*dtos.GlobalVar, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, errorwrap.WrapError("Unwilling to list global variables", errorwrap.UnauthorizedReadErr(err))
	}

	var globalVars = make([]models.GlobalVar, 0)
	err := db.Select(&globalVars, sq.Select("*").
		From("global_vars").
		OrderBy("name ASC"))

	if err != nil {
		return nil, errorwrap.WrapError("Cannot list global variables", errorwrap.DatabaseErr(err))
	}

	var globalVarsDTO = make([]*dtos.GlobalVar, len(globalVars))
	for i, globalVar := range globalVars {
		globalVarsDTO[i] = &dtos.GlobalVar{
			Name:  globalVar.Name,
			Value: globalVar.Value,
		}
	}

	return globalVarsDTO, nil
}

func UpdateGlobalVar(ctx context.Context, db *database.Connection, i UpdateGlobalVarInput) error {
	globalVar, err := LookupGlobalVar(db, i.Name)
	if err != nil {
		return errorwrap.WrapError("Unable to update operation", errorwrap.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.AdminUsersOnly{}); err != nil {
		return errorwrap.WrapError("Unwilling to update operation", errorwrap.UnauthorizedWriteErr(err))
	}

	var val string
	var name string

	if i.Value != "" {
		val = i.Value
	} else {
		val = globalVar.Value
	}

	formattedName := helpers.StrToUpperCaseUnderscore(i.NewName)

	if i.NewName != "" {
		name = formattedName
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
			return errorwrap.BadInputErr(errorwrap.WrapError("Global variable already exists", err), "A global variable with this name already exists")
		}
		return errorwrap.WrapError("Cannot update global variable", errorwrap.DatabaseErr(err))
	}

	return nil
}
