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

// DeleteAuthSchemeUsers removes/unlinks all users from a provided scheme
func DeleteAuthSchemeUsers(ctx context.Context, db *database.Connection, schemeCode string) error {
	if err := policy.Require(middleware.Policy(ctx), policy.CanDeleteAuthForAllUsers{SchemeCode: schemeCode}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	err := db.Delete(sq.Delete("auth_scheme_data").Where(sq.Eq{"auth_scheme": schemeCode}))
	if err != nil {
		return backend.DatabaseErr(err)
	}

	return nil
}
