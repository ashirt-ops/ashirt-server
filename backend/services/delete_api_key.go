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

type DeleteAPIKeyInput struct {
	AccessKey string
	UserSlug  string
}

func DeleteAPIKey(ctx context.Context, db *database.Connection, i DeleteAPIKeyInput) error {
	var userID int64
	var err error

	if userID, err = selfOrSlugToUserID(ctx, db, i.UserSlug); err != nil {
		return backend.DatabaseErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyAPIKeys{UserID: userID}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	var apiKeyID int64

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Get(&apiKeyID, sq.Select("id").
			From("api_keys").
			Where(sq.Eq{"user_id": userID, "access_key": i.AccessKey}))
		tx.Delete(sq.Delete("api_keys").Where(sq.Eq{"id": apiKeyID}))
	})
	if err != nil {
		if database.IsEmptyResultSetError(err) {
			return backend.UnauthorizedWriteErr(err)
		}
		return backend.DatabaseErr(err)
	}

	return nil
}
