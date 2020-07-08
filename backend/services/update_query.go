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

type UpdateQueryInput struct {
	OperationSlug string
	ID            int64
	Name          string
	Query         string
}

// UpdateQuery modifies a query for the given operation
func UpdateQuery(ctx context.Context, db *database.Connection, i UpdateQueryInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyQueriesOfOperation{OperationID: operation.ID}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	if i.Name == "" {
		return backend.MissingValueErr("Name")
	}
	if i.Query == "" {
		return backend.MissingValueErr("Query")
	}

	ub := sq.Update("queries").
		SetMap(map[string]interface{}{
			"name":  i.Name,
			"query": i.Query,
		}).
		Where(sq.Eq{"id": i.ID, "operation_id": operation.ID})

	err = db.Update(ub)
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return backend.BadInputErr(err, "A saved query with this name or query already exists")
		}
		return backend.UnauthorizedWriteErr(err)
	}
	return nil
}
