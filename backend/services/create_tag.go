// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
)

type CreateTagInput struct {
	Name          string
	ColorName     string
	OperationSlug string
}

func CreateTag(ctx context.Context, db *database.Connection, i CreateTagInput) (*dtos.Tag, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedWriteErr(err)
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}

	tagID, err := db.Insert("tags", map[string]interface{}{
		"name":         i.Name,
		"color_name":   i.ColorName,
		"operation_id": operation.ID,
	})
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}
	return &dtos.Tag{
		ID:        tagID,
		Name:      i.Name,
		ColorName: i.ColorName,
	}, nil
}
