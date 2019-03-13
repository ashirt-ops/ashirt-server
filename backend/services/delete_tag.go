// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type DeleteTagInput struct {
	ID            int64
	OperationSlug string
}

// DeleteTag removes a tag and untags all evidence with the tag
func DeleteTag(ctx context.Context, db *database.Connection, i DeleteTagInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"tag_id": i.ID}))
		tx.Delete(sq.Delete("tags").Where(sq.Eq{"id": i.ID}))
	})
	if err != nil {
		return backend.DatabaseErr(err)
	}

	return nil
}
