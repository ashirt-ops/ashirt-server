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

// QueriesForOperationOutput contains the given name of a query ("All evidence tagged with 'big deal'")
// as well as the actual, full length query ('tag: "big deal"')
type QueriesForOperationOutput struct {
	ID    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Query string `db:"query" json:"query"`
	Type  string `db:"type" json:"type"`
}

// ListQueriesForOperation retrieves all saved queries for a given operation id
func ListQueriesForOperation(ctx context.Context, db *database.Connection, operationSlug string) ([]QueriesForOperationOutput, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var queries = make([]QueriesForOperationOutput, 0)

	err = db.Select(&queries, sq.Select("id", "name", "query", "type").
		From("queries").
		Where(sq.Eq{"operation_id": operation.ID}).
		OrderBy("name ASC"))

	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	return queries, nil
}
