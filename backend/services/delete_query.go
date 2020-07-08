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

type DeleteQueryInput struct {
	OperationSlug string
	ID            int64
}

// DeleteQuery removes a saved query for the given operation
func DeleteQuery(ctx context.Context, db *database.Connection, i DeleteQueryInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyQueriesOfOperation{OperationID: operation.ID}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	err = db.Delete(sq.Delete("queries").Where(sq.Eq{"id": i.ID, "operation_id": operation.ID}))
	if err != nil {
		return backend.DatabaseErr(err)
	}

	return nil
}
