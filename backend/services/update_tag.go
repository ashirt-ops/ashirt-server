// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"

	sq "github.com/Masterminds/squirrel"
)

type UpdateDefaultTagInput struct {
	ID        int64
	Name      string
	ColorName string
}

type UpdateTagInput struct {
	ID            int64
	OperationSlug string
	Name          string
	ColorName     string
}

// UpdateTag updates a tag's name and color
func UpdateTag(ctx context.Context, db *database.Connection, i UpdateTagInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to update tag", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to update tag", backend.UnauthorizedWriteErr(err))
	}

	err = db.Update(sq.Update("tags").
		SetMap(map[string]interface{}{
			"name":       i.Name,
			"color_name": i.ColorName,
		}).
		Where(sq.Eq{"id": i.ID}))

	if err != nil {
		return backend.WrapError("Cannot update tag", backend.DatabaseErr(err))
	}
	return nil
}

func UpdateDefaultTag(ctx context.Context, db *database.Connection, i UpdateDefaultTagInput) error {
	if err := policyRequireWithAdminBypass(ctx, policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Unwilling to update default tag", backend.UnauthorizedWriteErr(err))
	}

	err := db.Update(sq.Update("default_tags").
		SetMap(map[string]interface{}{
			"name":       i.Name,
			"color_name": i.ColorName,
		}).
		Where(sq.Eq{"id": i.ID}))

	if err != nil {
		return backend.WrapError("Cannot update default tag", backend.DatabaseErr(err))
	}
	return nil
}
