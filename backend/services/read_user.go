// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
)

// ReadUser retrieves a detailed view of a user. This is separate from the data retriving by listing
// users, or reading another user's profile (when not an admin)
func ReadUser(ctx context.Context, db *database.Connection, userSlug string) (*dtos.UserOwnView, error) {
	userID, err := selfOrSlugToUserID(ctx, db, userSlug)
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadDetailedUser{UserID: userID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var user models.User
	var authSchemes []models.AuthSchemeData
	err = db.WithTx(ctx, func(tx *database.Transactable) {
		db.Get(&user, sq.Select("first_name", "last_name", "slug", "email", "admin", "headless").
			From("users").
			Where(sq.Eq{"id": userID}))

		db.Select(&authSchemes, sq.Select("user_key", "auth_scheme", "last_login").
			From("auth_scheme_data").
			Where(sq.Eq{"user_id": userID}))
	})
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	auths := make([]dtos.AuthenticationInfo, len(authSchemes))
	for i, v := range authSchemes {
		auths[i] = dtos.AuthenticationInfo{
			UserKey:        v.UserKey,
			AuthSchemeCode: v.AuthScheme,
			AuthLogin:      v.LastLogin,
		}
	}

	return &dtos.UserOwnView{
		User: dtos.User{
			Slug:      user.Slug,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		Email:          user.Email,
		Admin:          user.Admin,
		Headless:       user.Headless,
		Authentication: auths,
	}, nil
}
