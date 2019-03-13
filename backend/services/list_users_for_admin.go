// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"strings"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"

	sq "github.com/Masterminds/squirrel"
)

type ListUsersForAdminInput struct {
	UserFilter
	Pagination
	IncludeDeleted bool
}

// ListUsersForAdmin retreives standard User (public) details, and aguments with some particular fields
// meant for admin review. For use in admin views only.
func ListUsersForAdmin(ctx context.Context, db *database.Connection, i ListUsersForAdminInput) (*dtos.PaginationWrapper, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var users []struct {
		models.User
		AuthSchemes *string `db:"auth_schemes"`
	}

	sb := sq.Select("slug", "first_name", "last_name", "email", "admin", "disabled", "headless", "deleted_at", "GROUP_CONCAT(auth_scheme) AS auth_schemes").
		From("users").
		LeftJoin("auth_scheme_data ON auth_scheme_data.user_id = users.id").
		GroupBy("users.id")

	i.AddWhere(&sb)

	if !i.IncludeDeleted {
		sb = sb.Where(sq.Eq{"deleted_at": nil})
	}

	err := i.Pagination.Select(ctx, db, &users, sb)
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	usersDTO := []*dtos.UserAdminView{}
	for _, user := range users {
		// Group_Concat will return null if there are no authentication schemes listed. The below forces the schemes into a slice to avoid errors.
		authSchemes := []string{}
		if user.AuthSchemes != nil {
			authSchemes = strings.Split(*user.AuthSchemes, ",")
		}

		usersDTO = append(usersDTO, &dtos.UserAdminView{
			User: dtos.User{
				Slug:      user.Slug,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
			Email:       user.Email,
			Admin:       user.Admin,
			Headless:    user.Headless,
			AuthSchemes: authSchemes,
			Disabled:    user.Disabled,
			Deleted:     user.DeletedAt != nil,
		})
	}

	return i.Pagination.WrapData(usersDTO), nil
}
