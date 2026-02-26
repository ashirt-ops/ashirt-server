// Package services handles the logic behind all of the Web/API actions
package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/dtos"
	"github.com/ashirt-ops/ashirt-server/internal/errorwrap"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

const accessKeyLength = 18
const secretKeyLength = 64

type DeleteAPIKeyInput struct {
	AccessKey string
	UserSlug  string
}

func CreateAPIKey(ctx context.Context, db *database.Connection, userSlug string) (*dtos.APIKey, error) {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, userSlug); err != nil {
		return nil, errorwrap.WrapError("Unable to create api key", errorwrap.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyAPIKeys{UserID: userID}); err != nil {
		return nil, errorwrap.WrapError("Unable to create api key", errorwrap.UnauthorizedWriteErr(err))
	}

	accessKey := make([]byte, accessKeyLength)
	if _, err := rand.Read(accessKey); err != nil {
		return nil, errorwrap.WrapError("Unable to generate api key", err)
	}
	accessKeyStr := base64.URLEncoding.EncodeToString(accessKey)

	secretKey := make([]byte, secretKeyLength)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, errorwrap.WrapError("Unable to create secret key", err)
	}

	prefixedAccessKey := "AS-" + accessKeyStr

	_, err = db.Insert("api_keys", map[string]interface{}{
		"user_id":    userID,
		"access_key": prefixedAccessKey,
		"secret_key": secretKey,
	})
	if err != nil {
		return nil, errorwrap.WrapError("Unable to record api and secret keys", errorwrap.DatabaseErr(err))
	}

	return &dtos.APIKey{
		AccessKey: prefixedAccessKey,
		SecretKey: secretKey,
	}, nil
}

func DeleteAPIKey(ctx context.Context, db *database.Connection, i DeleteAPIKeyInput) error {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, i.UserSlug); err != nil {
		return errorwrap.WrapError("Unable to delete API Key", errorwrap.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyAPIKeys{UserID: userID}); err != nil {
		return errorwrap.WrapError("Unwilling to delete API Key", errorwrap.UnauthorizedWriteErr(err))
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
			return errorwrap.WrapError("API key does not exist", errorwrap.UnauthorizedWriteErr(err))
		}
		return errorwrap.WrapError("Cannot delete API key", errorwrap.DatabaseErr(err))
	}

	return nil
}

func ListAPIKeys(ctx context.Context, db *database.Connection, userSlug string) ([]*dtos.APIKey, error) {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, userSlug); err != nil {
		return nil, errorwrap.WrapError("Unable to list api keys", errorwrap.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanListAPIKeys{UserID: userID}); err != nil {
		return nil, errorwrap.WrapError("Unwilling to list api keys", errorwrap.UnauthorizedReadErr(err))
	}

	var keys []models.APIKey
	err = db.Select(&keys, sq.Select("access_key", "last_auth").
		From("api_keys").
		Where(sq.Eq{"user_id": userID}))

	if err != nil {
		return nil, errorwrap.WrapError("Cannot list api keys", errorwrap.DatabaseErr(err))
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
