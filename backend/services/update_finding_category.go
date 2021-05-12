// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
)

type UpdateFindingCategoryInput struct {
	ID       int64
	Category string
}

// UpdateFindingCategory updates the specified entry in the finding_categories table
func UpdateFindingCategory(ctx context.Context, db *database.Connection, i UpdateFindingCategoryInput) error {
	if err := isAdmin(ctx); err != nil {
		return backend.WrapError("Unable to update the finding category", backend.UnauthorizedWriteErr(err))
	}

	err := db.Update(sq.Update("finding_categories").
		SetMap(map[string]interface{}{
			"category": i.Category,
		}).
		Where(sq.Eq{"id": i.ID}))

	if err != nil {
		return backend.WrapError("Cannot update finding category", backend.DatabaseErr(err))
	}
	return nil
}
