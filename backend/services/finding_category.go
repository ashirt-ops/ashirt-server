// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"

	sq "github.com/Masterminds/squirrel"
)

type DeleteFindingCategoryInput struct {
	FindingCategoryID int64
	DoDelete          bool
}

type UpdateFindingCategoryInput struct {
	ID       int64
	Category string
}

// CreateFindingCategory adds a new finding category to the finding_categories table
func CreateFindingCategory(ctx context.Context, db *database.Connection, newCategory string) (*dtos.FindingCategory, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unable to create a new finding category", backend.UnauthorizedWriteErr(err))
	}

	id, err := db.Insert("finding_categories", map[string]interface{}{
		"category": newCategory,
	})
	if err != nil {
		return nil, err
	}

	return &dtos.FindingCategory{ID: id, Category: newCategory}, nil
}

// DeleteFindingCategory removes an entry from the finding_categories table
func DeleteFindingCategory(ctx context.Context, db *database.Connection, i DeleteFindingCategoryInput) error {
	if err := isAdmin(ctx); err != nil {
		return backend.WrapError("Unable to delete a finding category", backend.UnauthorizedWriteErr(err))
	}

	query := sq.Update("finding_categories").
		Where(sq.Eq{"id": i.FindingCategoryID})

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
