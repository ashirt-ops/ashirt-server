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

type DeleteAuthSchemeInput struct {
	UserSlug   string
	SchemeName string
}

// DeleteAuthScheme removes a user's association with a particular auth_scheme. This function applies for both
// admin related actions and plain user actions. If UserSlug is not provided, this will apply to the requesting
// user. If it is provided, then this triggers admin validation, and will apply to the provided user matching the
// given slug.
func DeleteAuthScheme(ctx context.Context, db *database.Connection, i DeleteAuthSchemeInput) error {
	var userID int64
	var err error

	if userID, err = selfOrSlugToUserID(ctx, db, i.UserSlug); err != nil {
		return backend.DatabaseErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanDeleteAuthScheme{UserID: userID, SchemeCode: i.SchemeName}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	err = db.Delete(sq.Delete("auth_scheme_data").Where(sq.Eq{"user_id": userID, "auth_scheme": i.SchemeName}))
	if err != nil {
		return backend.DatabaseErr(err)
	}

	return nil
}
