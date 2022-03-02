// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type DeleteTagInput struct {
	ID            int64
	OperationSlug string
}

type DeleteDefaultTagInput struct {
	ID int64
}

// DeleteTag removes a tag and untags all evidence with the tag
func DeleteTag(ctx context.Context, db *database.Connection, i DeleteTagInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to delete tag", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete tag", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"tag_id": i.ID}))
		tx.Delete(sq.Delete("tags").Where(sq.Eq{"id": i.ID}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete tag", backend.DatabaseErr(err))
	}

	return nil
}

// DeleteDefaultTag removes a single tag in the default_tags table by the tag id. Admin only.
func DeleteDefaultTag(ctx context.Context, db *database.Connection, i DeleteDefaultTagInput) error {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Unwilling to delete default tag", backend.UnauthorizedWriteErr(err))
	}

	err := db.Delete(sq.Delete("default_tags").Where(sq.Eq{"id": i.ID}))
	if err != nil {
		return backend.WrapError("Cannot delete default tag", backend.DatabaseErr(err))
	}

	return nil
}
