package services

import (
	"context"
	"fmt"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/dtos"
	"github.com/ashirt-ops/ashirt-server/internal/errorwrap"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

// CreateQueryInput provides a structure that holds the values needed to generate a new saved query
type CreateQueryInput struct {
	OperationSlug string
	Name          string
	Query         string
	Type          string
}

type DeleteQueryInput struct {
	OperationSlug string
	ID            int64
}

type UpdateQueryInput struct {
	OperationSlug string
	ID            int64
	Name          string
	Query         string
}

type UpsertQueryInput struct {
	CreateQueryInput
	ReplaceName bool
}

// CreateQuery inserts a new query into the database
func CreateQuery(ctx context.Context, db *database.Connection, i CreateQueryInput) (*dtos.Query, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, errorwrap.WrapError("Unable to create query", errorwrap.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyQueriesOfOperation{OperationID: operation.ID}); err != nil {
		return nil, errorwrap.WrapError("Unable to create query", errorwrap.UnauthorizedWriteErr(err))
	}

	validationError := validateCreateQueryInput(i)
	if validationError != nil {
		return nil, errorwrap.WrapError("CreateQuery validation  error", validationError)
	}

	queryID, err := db.Insert("queries", map[string]interface{}{
		"operation_id": operation.ID,
		"name":         i.Name,
		"query":        i.Query,
		"type":         i.Type,
	})
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return nil, errorwrap.BadInputErr(errorwrap.WrapError("Query already exists", err), "A query with this name already exists")
		}
		return nil, errorwrap.WrapError("Unable to add new query", errorwrap.DatabaseErr(err))
	}

	return &dtos.Query{
		ID:    queryID,
		Name:  i.Name,
		Query: i.Query,
		Type:  i.Type,
	}, nil
}

// DeleteQuery removes a saved query for the given operation
func DeleteQuery(ctx context.Context, db *database.Connection, i DeleteQueryInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return errorwrap.WrapError("Unable to delete query", errorwrap.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyQueriesOfOperation{OperationID: operation.ID}); err != nil {
		return errorwrap.WrapError("Unwilling to delete query", errorwrap.UnauthorizedWriteErr(err))
	}

	err = db.Delete(sq.Delete("queries").Where(sq.Eq{"id": i.ID, "operation_id": operation.ID}))
	if err != nil {
		return errorwrap.WrapError("Cannot delete query", errorwrap.DatabaseErr(err))
	}

	return nil
}

// ListQueriesForOperation retrieves all saved queries for a given operation id
func ListQueriesForOperation(ctx context.Context, db *database.Connection, operationSlug string) ([]*dtos.Query, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, errorwrap.WrapError("Unable to list queries", errorwrap.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, errorwrap.WrapError("Unwilling to list queries", errorwrap.UnauthorizedReadErr(err))
	}

	var queries = make([]models.Query, 0)
	err = db.Select(&queries, sq.Select("id", "name", "query", "type").
		From("queries").
		Where(sq.Eq{"operation_id": operation.ID}).
		OrderBy("name ASC"))

	if err != nil {
		return nil, errorwrap.WrapError("Cannot list queries", errorwrap.DatabaseErr(err))
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

// UpdateQuery modifies a query for the given operation
func UpdateQuery(ctx context.Context, db *database.Connection, i UpdateQueryInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return errorwrap.WrapError("Unable to update query", errorwrap.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyQueriesOfOperation{OperationID: operation.ID}); err != nil {
		return errorwrap.WrapError("Unwilling to update query", errorwrap.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return errorwrap.MissingValueErr("Name")
	}
	if i.Query == "" {
		return errorwrap.MissingValueErr("Query")
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
			return errorwrap.WrapError("Cannot update query", errorwrap.BadInputErr(err, "A saved query with this name or query already exists"))
		}
		return errorwrap.WrapError("Cannot update query", errorwrap.UnauthorizedWriteErr(err))
	}
	return nil
}

func UpsertQuery(ctx context.Context, db *database.Connection, i UpsertQueryInput) (*dtos.Query, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, errorwrap.WrapError("Unable to upsert query", errorwrap.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyQueriesOfOperation{OperationID: operation.ID}); err != nil {
		return nil, errorwrap.WrapError("Unable to upsert query", errorwrap.UnauthorizedWriteErr(err))
	}

	validationError := validateCreateQueryInput(i.CreateQueryInput)
	if validationError != nil {
		return nil, errorwrap.WrapError("UpsertQuery validation error", validationError)
	}

	onDuplicates := "ON DUPLICATE KEY UPDATE "

	if i.ReplaceName {
		onDuplicates += "name=VALUES(name)"
	} else {
		onDuplicates += "query=VALUES(query)"
	}

	queryID, err := db.Insert("queries", map[string]interface{}{
		"operation_id": operation.ID,
		"name":         i.Name,
		"query":        i.Query,
		"type":         i.Type,
	}, onDuplicates)
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return nil, errorwrap.BadInputErr(errorwrap.WrapError("Query already exists", err), "A query with this name already exists")
		}
		return nil, errorwrap.WrapError("Unable to add new query", errorwrap.DatabaseErr(err))
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
		return errorwrap.MissingValueErr("Query")
	}
	if input.Name == "" {
		return errorwrap.MissingValueErr("Name")
	}
	if input.Type != "findings" && input.Type != "evidence" {
		err := fmt.Errorf("Bad type: %s", input.Type)
		return errorwrap.BadInputErr(err, err.Error())
	}
	return nil
}
