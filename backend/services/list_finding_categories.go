// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"

	sq "github.com/Masterminds/squirrel"
)

// ListFindingCategories retrieves a list of all of the finding categories present in the database.
func ListFindingCategories(ctx context.Context, db *database.Connection, includeDeleted bool) (interface{}, error) {
	query := sq.Select("fc.id", "fc.category", "fc.deleted_at", "count(f.id) AS usage_count").
		From("finding_categories AS fc").
		LeftJoin("findings AS f ON fc.id = f.category_id").
		GroupBy("fc.category").
		OrderBy("fc.category")

	if !includeDeleted {
		query = query.Where(sq.Eq{"deleted_at": nil})
	}

	type FindingCategoryWithUsage struct {
		ID         int64      `db:"id"`
		Category   string     `db:"category"`
		UsageCount int64      `db:"usage_count"`
		DeletedAt  *time.Time `db:"deleted_at"`
	}

	var categories []FindingCategoryWithUsage

	if err := db.Select(&categories, query); err != nil {
		return nil, backend.WrapError("Cannot list finding categories", backend.DatabaseErr(err))
	}

	rtn := make([]*dtos.FindingCategory, len(categories))
	for i, cat := range categories {
		rtn[i] = &dtos.FindingCategory{
			ID:         cat.ID,
			Category:   cat.Category,
			Deleted:    cat.DeletedAt != nil,
			UsageCount: cat.UsageCount,
		}
	}

	return rtn, nil
}
