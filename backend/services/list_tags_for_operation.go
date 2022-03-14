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
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ListTagsForOperationInput struct {
	OperationSlug string
}

func ListTagsForOperation(ctx context.Context, db *database.Connection, i ListTagsForOperationInput) ([]*dtos.TagWithUsage, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list tags for operation", backend.UnauthorizedReadErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list tags for operation", backend.UnauthorizedReadErr(err))
	}

	return listTagsForOperation(db, operation.ID)
}

// listTagsForOperation generates a list tags associted with a given operation. This does not
// check permission, and so is not exported, and is intended to only be used as a helper method
// for other services
func listTagsForOperation(db *database.Connection, operationID int64) ([]*dtos.TagWithUsage, error) {
	type DBTag struct {
		models.Tag
		TagCount int64 `db:"tag_count"`
	}
	var tags []DBTag
	err := db.Select(&tags, sq.Select("tags.*").Column("count(tag_id) AS tag_count").
		From("tags").
		LeftJoin("tag_evidence_map ON tag_evidence_map.tag_id = tags.id").
		Where(sq.Eq{"operation_id": operationID}).
		GroupBy("tags.id").
		OrderBy("tags.id ASC"))
	if err != nil {
		return nil, backend.WrapError("Cannot get tags for operation", backend.DatabaseErr(err))
	}

	tagsDTO := make([]*dtos.TagWithUsage, len(tags))
	for idx, tag := range tags {
		tagsDTO[idx] = &dtos.TagWithUsage{
			Tag: dtos.Tag{
				ID:        tag.Tag.ID,
				Name:      tag.Tag.Name,
				ColorName: tag.Tag.ColorName,
			},
			EvidenceCount: tag.TagCount,
		}
	}
	return tagsDTO, nil
}

// ListDefaultTags provides a list of all of the tags in the default_tags table. Admin only.
func ListDefaultTags(ctx context.Context, db *database.Connection) ([]*dtos.DefaultTag, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.WrapError("Unwilling to list default tags", backend.UnauthorizedReadErr(err))
	}

	var tags []models.Tag
	err := db.Select(&tags, sq.Select("id", "name", "color_name").From("default_tags"))

	if err != nil {
		return nil, backend.WrapError("Cannot get default tags", backend.DatabaseErr(err))
	}

	tagsDTO := make([]*dtos.DefaultTag, len(tags))
	for idx, tag := range tags {
		tagsDTO[idx] = &dtos.DefaultTag{
			ID:        tag.ID,
			Name:      tag.Name,
			ColorName: tag.ColorName,
		}
	}
	return tagsDTO, nil
}
