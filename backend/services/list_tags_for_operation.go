// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"

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

	var tags []models.Tag
	err = db.Select(&tags, sq.Select("id", "name", "color_name").
		From("tags").
		Where(sq.Eq{"operation_id": operation.ID}).
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
