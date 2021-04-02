// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend/dtos"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
)

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

	return &dtos.FindingCategory{ ID: id, Category: newCategory }, nil
}
