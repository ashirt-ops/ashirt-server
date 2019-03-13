// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
)

// CreateQueryInput provides a structure that holds the values needed to generate a new saved query
type CreateQueryInput struct {
	OperationSlug string
	Name          string
	Query         string
	Type          string
}

// CreateQuery inserts a new query into the database
func CreateQuery(ctx context.Context, db *database.Connection, i CreateQueryInput) (*dtos.Query, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.UnauthorizedWriteErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyQueriesOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedWriteErr(err)
	}

	validationError := validateCreateQueryInput(i)
	if validationError != nil {
		return nil, validationError
	}

	queryID, err := db.Insert("queries", map[string]interface{}{
		"operation_id": operation.ID,
		"name":         i.Name,
		"query":        i.Query,
		"type":         i.Type,
	})
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return nil, backend.BadInputErr(err, "A query with this name already exists")
		}
		return nil, backend.DatabaseErr(err)
	}

	return &dtos.Query{
		ID:    queryID,
		Name:  i.Name,
		Query: i.Query,
		Type:  i.Type,
	}, nil
}

func validateCreateQueryInput(input CreateQueryInput) error {
	if input.Query == "" {
		return backend.MissingValueErr("Query")
	}
	if input.Name == "" {
		return backend.MissingValueErr("Name")
	}
	if input.Type != "findings" && input.Type != "evidence" {
		err := fmt.Errorf("Bad type: %s", input.Type)
		return backend.BadInputErr(err, err.Error())
	}
	return nil
}
