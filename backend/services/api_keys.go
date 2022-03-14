// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type DeleteAPIKeyInput struct {
	AccessKey string
	UserSlug  string
}

type RotateAPIKeyInput = DeleteAPIKeyInput

const accessKeyLength = 18
const secretKeyLength = 64

func CreateAPIKey(ctx context.Context, db *database.Connection, userSlug string) (*dtos.APIKey, error) {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, userSlug); err != nil {
		return nil, backend.WrapError("Unable to create api key", backend.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyAPIKeys{UserID: userID}); err != nil {
		return nil, backend.WrapError("Unable to create api key", backend.UnauthorizedWriteErr(err))
	}

	return createAPIKey(db, userID)
}

func DeleteAPIKey(ctx context.Context, db *database.Connection, i DeleteAPIKeyInput) error {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, i.UserSlug); err != nil {
		return backend.WrapError("Unable to delete API Key", backend.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyAPIKeys{UserID: userID}); err != nil {
		return backend.WrapError("Unwilling to delete API Key", backend.UnauthorizedWriteErr(err))
	}

	return deleteAPIKey(ctx, db, userID, i.AccessKey)
}

func ListAPIKeys(ctx context.Context, db *database.Connection, userSlug string) ([]*dtos.APIKey, error) {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, userSlug); err != nil {
		return nil, backend.WrapError("Unable to list api keys", backend.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanListAPIKeys{UserID: userID}); err != nil {
		return nil, backend.WrapError("Unwilling to list api keys", backend.UnauthorizedReadErr(err))
	}

	var keys []models.APIKey
	err = db.Select(&keys, sq.Select("access_key", "last_auth").
		From("api_keys").
		Where(sq.Eq{"user_id": userID}))

	if err != nil {
		return nil, backend.WrapError("Cannot list api keys", backend.DatabaseErr(err))
	}

	keysDTO := make([]*dtos.APIKey, len(keys))
	for i, key := range keys {
		keysDTO[i] = &dtos.APIKey{
			AccessKey: key.AccessKey,
			LastAuth:  key.LastAuth,
		}
	}
	return keysDTO, nil
}

func createAPIKey(db database.ConnectionProxy, userID int64) (*dtos.APIKey, error) {
	accessKey := make([]byte, accessKeyLength)
	if _, err := rand.Read(accessKey); err != nil {
		return nil, backend.WrapError("Unable to generate api key", err)
	}
	accessKeyStr := base64.URLEncoding.EncodeToString(accessKey)

	secretKey := make([]byte, secretKeyLength)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, backend.WrapError("Unable to create secret key", err)
	}

	_, err := db.Insert("api_keys", map[string]interface{}{
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

func deleteAPIKey(ctx context.Context, db database.ConnectionProxy, userID int64, accessKey string) error {
	var apiKeyID int64

	err := db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Get(&apiKeyID, sq.Select("id").
			From("api_keys").
			Where(sq.Eq{"user_id": userID, "access_key": accessKey}))
		tx.Delete(sq.Delete("api_keys").Where(sq.Eq{"id": apiKeyID}))
	})
	if err != nil {
		if database.IsEmptyResultSetError(err) {
			return backend.WrapError("API key does not exist", backend.UnauthorizedWriteErr(err))
		}
		return backend.WrapError("Cannot delete API key", backend.DatabaseErr(err))
	}

	return nil
}

func rotateAPIKey(ctx context.Context, db database.ConnectionProxy, userID int64, accessKey string) (*dtos.APIKey, error) {
	var apiKey *dtos.APIKey

	err := db.WithTx(ctx, func(tx *database.Transactable) {
		err := deleteAPIKey(ctx, tx, userID, accessKey)
		if err != nil {
			tx.FailTransaction(err)
		}
		apiKey, err = createAPIKey(tx, userID)
		if err != nil {
			tx.FailTransaction(err)
		}
	})

	return apiKey, err
}

func RotateAPIKey(ctx context.Context, db *database.Connection, i RotateAPIKeyInput) (*dtos.APIKey, error) {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, i.UserSlug); err != nil {
		return nil, backend.WrapError("Unable to delete API Key", backend.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyAPIKeys{UserID: userID}); err != nil {
		return nil, backend.WrapError("Unwilling to delete API Key", backend.UnauthorizedWriteErr(err))
	}

	return rotateAPIKey(ctx, db, userID, i.AccessKey)
}
