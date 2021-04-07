// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"

	sq "github.com/Masterminds/squirrel"
)

// ListFindingCategories retrieves a list of all of the finding categories present in the database.
func ListFindingCategories(ctx context.Context, db *database.Connection, includeDeleted bool) (interface{}, error) {
	query := sq.Select("id", "category", "deleted_at").
		From("finding_categories").
		OrderBy("category")

	if !includeDeleted {
		query = query.Where(sq.Eq{"deleted_at": nil})
	}

	var categories []models.FindingCategory
	if err := db.Select(&categories, query); err != nil {
		return nil, backend.WrapError("Cannot list finding categories", backend.DatabaseErr(err))
	}

	rtn := make([]*dtos.FindingCategory, len(categories))
	for i, cat := range categories {
		rtn[i] = &dtos.FindingCategory{
			ID:       cat.ID,
			Category: cat.Category,
			Deleted:  cat.DeletedAt != nil,
		}
	}

	return rtn, nil
}
