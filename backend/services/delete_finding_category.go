// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
)

// DeleteFindingCategory removes an entry from the finding_categories table
func DeleteFindingCategory(ctx context.Context, db *database.Connection, findingCategoryId int64) error {
	if err := isAdmin(ctx); err != nil {
		return backend.WrapError("Unable to delete a finding category", backend.UnauthorizedWriteErr(err))
	}

	err := db.Delete(sq.Delete("finding_categories").
		Where(sq.Eq{"id": findingCategoryId}),
	)
	if err != nil {
		return backend.WrapError("Cannot delete finding category", backend.DatabaseErr(err))
	}

	return nil
}
