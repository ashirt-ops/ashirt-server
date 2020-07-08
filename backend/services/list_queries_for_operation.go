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

// ListQueriesForOperation retrieves all saved queries for a given operation id
func ListQueriesForOperation(ctx context.Context, db *database.Connection, operationSlug string) ([]*dtos.Query, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var queries = make([]models.Query, 0)
	err = db.Select(&queries, sq.Select("id", "name", "query", "type").
		From("queries").
		Where(sq.Eq{"operation_id": operation.ID}).
		OrderBy("name ASC"))

	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	var queriesDTO = make([]*dtos.Query, len(queries))
	for i, query := range queries {
		queriesDTO[i] = &dtos.Query{
			ID:    query.ID,
			Name:  query.Name,
			Query: query.Query,
			Type:  query.Type,
		}
	}

	return queriesDTO, nil
}
