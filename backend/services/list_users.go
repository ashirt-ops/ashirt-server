// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"strings"
	"unicode"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ListUsersInput struct {
	Query          string
	IncludeDeleted bool
}

func ListUsers(ctx context.Context, db *database.Connection, i ListUsersInput) ([]*dtos.User, error) {
	if strings.ContainsAny(i.Query, "%_") || strings.TrimFunc(i.Query, unicode.IsSpace) == "" {
		return []*dtos.User{}, nil
	}

	var users []models.User
	query := sq.Select("slug", "first_name", "last_name").
		From("users").
		Where(sq.Like{"concat(first_name, ' ', last_name)": "%" + strings.ReplaceAll(i.Query, " ", "%") + "%"}).
		OrderBy("first_name").
		Limit(10)
	if !i.IncludeDeleted {
		query = query.Where(sq.Eq{"deleted_at": nil})
	}
	err := db.Select(&users, query)
	if err != nil {
		return nil, backend.WrapError("Cannot list users", backend.DatabaseErr(err))
	}

	usersDTO := []*dtos.User{}
	for _, user := range users {
		if middleware.Policy(ctx).Check(policy.CanReadUser{UserID: user.ID}) {
			usersDTO = append(usersDTO, &dtos.User{
				Slug:      user.Slug,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			})
		}
	}
	return usersDTO, nil
}
