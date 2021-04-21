// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ListEvidenceCreatorsForOperationInput struct {
	OperationSlug string
}

// ListEvidenceCreatorsForOperation returns a list of all users that have (ever) created a piece of
// evidence for a given operation slug. Note that this won't return users that _had_ created evidence
// that has since been deleted
func ListEvidenceCreatorsForOperation(ctx context.Context, db *database.Connection, i ListEvidenceCreatorsForOperationInput) ([]*dtos.User, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list evidence for an operation", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list evidence for an operation", backend.UnauthorizedReadErr(err))
	}

	var users []struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Slug      string `db:"slug"`
	}

	sb := sq.Select("users.slug", "users.first_name", "users.last_name").
		Distinct().
		From("operations").
		LeftJoin("evidence ON operations.id = evidence.operation_id").
		LeftJoin("users ON evidence.operator_id = users.id").
		Where(sq.Eq{"operations.slug": i.OperationSlug}).
		OrderBy("users.first_name ASC")

	err = db.Select(&users, sb)
	if err != nil {
		return nil, backend.WrapError("Cannot list evidence creators for an operation", backend.DatabaseErr(err))
	}

	usersDTO := make([]*dtos.User, len(users))
	for idx, user := range users {
		usersDTO[idx] = &dtos.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Slug:      user.Slug,
		}
	}
	return usersDTO, nil
}
