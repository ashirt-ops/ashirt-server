// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"

	sq "github.com/Masterminds/squirrel"
)

type ListTagsForOperationInput struct {
	OperationSlug string
}

func ListTagsForOperation(ctx context.Context, db *database.Connection, i ListTagsForOperationInput) ([]*dtos.Tag, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	return listTagsForOperation(db, operation.ID)
}

// listTagsForOperation generates a list tags associted with a given operation. This does not
// check permission, and so is not exported, and is intended to only be used as a helper method
// for other services
func listTagsForOperation(db *database.Connection, operationID int64) ([]*dtos.Tag, error) {
	var tags []models.Tag
	err := db.Select(&tags, sq.Select("id", "name", "color_name").
		From("tags").
		Where(sq.Eq{"operation_id": operationID}).
		OrderBy("id ASC"))
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	tagsDTO := make([]*dtos.Tag, len(tags))
	for idx, tag := range tags {
		tagsDTO[idx] = &dtos.Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			ColorName: tag.ColorName,
		}
	}
	return tagsDTO, nil
}
