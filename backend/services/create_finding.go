// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
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
		return nil, backend.UnauthorizedWriteErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedWriteErr(err)
	}

	if i.Title == "" {
		return nil, backend.MissingValueErr("Title")
	}

	if i.Category == "" {
		return nil, backend.MissingValueErr("Category")
	}

	findingUUID := uuid.New().String()
	_, err = db.Insert("findings", map[string]interface{}{
		"uuid":         findingUUID,
		"operation_id": operation.ID,
		"category":     i.Category,
		"title":        i.Title,
		"description":  i.Description,
	})
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	return &dtos.Finding{
		UUID:        findingUUID,
		Title:       i.Title,
		Description: i.Description,
	}, nil
}
