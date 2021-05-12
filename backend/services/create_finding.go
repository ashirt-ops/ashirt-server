// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
)

type CreateFindingInput struct {
	OperationSlug string
	Category      string
	Title         string
	Description   string
}

func CreateFinding(ctx context.Context, db *database.Connection, i CreateFindingInput) (*dtos.Finding, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to create finding", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unable to create finding", backend.UnauthorizedWriteErr(err))
	}

	if i.Title == "" {
		return nil, backend.MissingValueErr("Title")
	}

	if i.Category == "" {
		return nil, backend.MissingValueErr("Category")
	}

	useCategoryID, err := getFindingCategoryID(i.Category, db.Select)

	if err != nil {
		return nil, backend.WrapError("Unable create finding", err)
	}
	if useCategoryID == nil {
		return nil, backend.BadInputErr(errors.New("no such category"), "Unknown Category")
	}

	findingUUID := uuid.New().String()
	_, err = db.Insert("findings", map[string]interface{}{
		"uuid":         findingUUID,
		"operation_id": operation.ID,
		"category_id":  useCategoryID,
		"title":        i.Title,
		"description":  i.Description,
	})
	if err != nil {
		return nil, backend.WrapError("Unable to insert finding", backend.DatabaseErr(err))
	}

	return &dtos.Finding{
		UUID:        findingUUID,
		Title:       i.Title,
		Description: i.Description,
	}, nil
}
