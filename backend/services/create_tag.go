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

func MergeDefaultTags(ctx context.Context, db *database.Connection, i []CreateDefaultTagInput) error {
	if err := policyRequireWithAdminBypass(ctx, policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Unwilling to update default tag", backend.UnauthorizedWriteErr(err))
	}

	tagsToInsert := make([]CreateDefaultTagInput, 0, len(i))
	currentTagNames := make([]string, 0, len(i))

	for _, t := range i {
		if listContainsString(currentTagNames, t.Name) != -1 || t.Name == "" {
			continue // no need to re-process a tag if we've dealt with it -- just use the first instance
		} else {
			currentTagNames = append(currentTagNames, t.Name)
		}

		tagsToInsert = append(tagsToInsert, t)
	}

	err := db.BatchInsert("default_tags", len(tagsToInsert), func(idx int) map[string]interface{} {
		return map[string]interface{}{
			"name":       tagsToInsert[idx].Name,
			"color_name": tagsToInsert[idx].ColorName,
		}
	}, "ON DUPLICATE KEY UPDATE color_name=VALUES(color_name)")

	if err != nil {
		return backend.WrapError("Cannot update default tag", backend.DatabaseErr(err))
	}
	return nil
}

func listContainsString(haystack []string, needle string) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}
