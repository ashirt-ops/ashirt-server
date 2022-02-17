// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
)

type CreateTagInput struct {
	Name          string
	ColorName     string
	OperationSlug string
}

type CreateDefaultTagInput struct {
	Name      string
	ColorName string
}

func CreateTag(ctx context.Context, db *database.Connection, i CreateTagInput) (*dtos.Tag, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to create tag", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unable to create tag", backend.UnauthorizedWriteErr(err))
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
		return nil, backend.WrapError("Cannot add new tag", backend.DatabaseErr(err))
	}
	return &dtos.Tag{
		ID:        tagID,
		Name:      i.Name,
		ColorName: i.ColorName,
	}, nil
}

// CreateDefaultTag creates a single tag in the default_tags table. Admin only.
func CreateDefaultTag(ctx context.Context, db *database.Connection, i CreateDefaultTagInput) (*dtos.DefaultTag, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.WrapError("Unable to create default tag", backend.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}

	tagID, err := db.Insert("default_tags", map[string]interface{}{
		"name":       i.Name,
		"color_name": i.ColorName,
	})
	if err != nil {
		return nil, backend.WrapError("Cannot add new tag", backend.DatabaseErr(err))
	}
	return &dtos.DefaultTag{
		ID:        tagID,
		Name:      i.Name,
		ColorName: i.ColorName,
	}, nil
}
