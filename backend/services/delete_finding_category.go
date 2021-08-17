// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
)

type DeleteFindingCategoryInput struct {
	FindingCategoryId int64
	DoDelete          bool
}

// DeleteFindingCategory removes an entry from the finding_categories table
func DeleteFindingCategory(ctx context.Context, db *database.Connection, i DeleteFindingCategoryInput) error {
	if err := isAdmin(ctx); err != nil {
		return backend.WrapError("Unable to delete a finding category", backend.UnauthorizedWriteErr(err))
	}

	query := sq.Update("finding_categories").
		Where(sq.Eq{"id": i.FindingCategoryId})

	if i.DoDelete {
		query = query.Set("deleted_at", time.Now())
	} else {
		query = query.Set("deleted_at", nil)
	}

	if err := db.Update(query); err != nil {
		return backend.WrapError("Cannot delete finding category", backend.DatabaseErr(err))
	}

	return nil
}
