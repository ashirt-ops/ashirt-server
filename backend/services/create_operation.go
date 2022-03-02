// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateOperationInput struct {
	Slug    string
	OwnerID int64
	Name    string
}

func CreateOperation(ctx context.Context, db *database.Connection, i CreateOperationInput) (*dtos.Operation, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.CanCreateOperations{}); err != nil {
		return nil, backend.WrapError("Unable to create operation", backend.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}

	if i.Slug == "" {
		return nil, backend.MissingValueErr("Slug")
	}

	cleanSlug := SanitizeOperationSlug(i.Slug)
	if cleanSlug == "" {
		return nil, backend.BadInputErr(errors.New("Unable to create operation. Invalid operation slug"), "Slug must contain english letters or numbers")
	}

	err := db.WithTx(ctx, func(tx *database.Transactable) {
		operationID, _ := tx.Insert("operations", map[string]interface{}{
			"name":   i.Name,
			"status": models.OperationStatusPlanning,
			"slug":   cleanSlug,
		})
		tx.Insert("user_operation_permissions", map[string]interface{}{
			"user_id":      i.OwnerID,
			"operation_id": operationID,
			"role":         policy.OperationRoleAdmin,
		})

		// Copy default tags into new operation
		tx.Exec(sq.Insert("tags").
			Columns(
				"name", "color_name",
				"operation_id",
			).
			Select(sq.Select(
				"name", "color_name",
				fmt.Sprintf("%v AS operation_id", operationID),
			).From("default_tags")),
		)
	})
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return nil, backend.WrapError("Unable to create operation. Operation slug already exists.", backend.BadInputErr(err, "An operation with this slug already exists"))
		}
		return nil, backend.WrapError("Unable to add new operation", backend.DatabaseErr(err))
	}

	return &dtos.Operation{
		Slug:     cleanSlug,
		Name:     i.Name,
		NumUsers: 1,
		Status:   models.OperationStatusPlanning,
	}, nil
}

var disallowedCharactersRegex = regexp.MustCompile(`[^A-Za-z0-9]+`)

// SanitizeOperationSlug removes objectionable characters from a slug and returns the new slug.
// Current logic: only allow alphanumeric characters and hyphen, with hypen excluded at the start
// and end
func SanitizeOperationSlug(slug string) string {
	return strings.Trim(
		disallowedCharactersRegex.ReplaceAllString(strings.ToLower(slug), "-"),
		"-",
	)
}
