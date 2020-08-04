// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
)

const accessKeyLength = 18
const secretKeyLength = 64

func CreateAPIKey(ctx context.Context, db *database.Connection, userSlug string) (*dtos.APIKey, error) {
	var userID int64
	var err error

	if userID, err = selfOrSlugToUserID(ctx, db, userSlug); err != nil {
		return nil, backend.WrapError("Unable to create api key", backend.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyAPIKeys{UserID: userID}); err != nil {
		return nil, backend.WrapError("Unable to create api key", backend.UnauthorizedWriteErr(err))
	}

	accessKey := make([]byte, accessKeyLength)
	if _, err := rand.Read(accessKey); err != nil {
		return nil, backend.WrapError("Unable to generate api key", err)
	}
	accessKeyStr := base64.URLEncoding.EncodeToString(accessKey)

	secretKey := make([]byte, secretKeyLength)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, backend.WrapError("Unable to create secret key", err)
	}

	_, err = db.Insert("api_keys", map[string]interface{}{
		"user_id":    userID,
		"access_key": accessKeyStr,
		"secret_key": secretKey,
	})
	if err != nil {
		return nil, backend.WrapError("Unable to record api and secret keys", backend.DatabaseErr(err))
	}

	return &dtos.APIKey{
		AccessKey: accessKeyStr,
		SecretKey: secretKey,
	}, nil
}
