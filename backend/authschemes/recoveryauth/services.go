// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package recoveryauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

// deleteExpiredRecoveryCodes removes recovery codes from the database that are older than X minutes old
// Duration is determined by looking at the environment configuration
func deleteExpiredRecoveryCodes(ctx context.Context, db *database.Connection, expiryInMinutes int64) error {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	err := db.Delete(sq.Delete("auth_scheme_data").
		Where(sq.Eq{"auth_scheme": recoveryConsts.Code}).
		Where("TIMESTAMPDIFF(minute, created_at, ?) >= ?", time.Now(), expiryInMinutes)) // ensure timestamps are sufficently old
	if err != nil {
		return backend.DatabaseErr(err)
	}

	return nil
}

// generateRecoveryCodeForUser creates a new, cryptographically random, base64-encoded,
// recovery key. This key is then attached to a user as a authorization method.
func generateRecoveryCodeForUser(ctx context.Context, bridge authschemes.AShirtAuthBridge, userSlug string) (interface{}, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.UnauthorizedWriteErr(err)
	}

	userID, err := bridge.GetUserIDFromSlug(userSlug)
	if err != nil {
		return nil, err
	}

	authKey := make([]byte, recoveryKeyLength)
	if _, err := rand.Read(authKey); err != nil {
		return nil, err
	}
	authKeyStr := base64.URLEncoding.EncodeToString(authKey)

	err = bridge.CreateNewAuthForUser(authschemes.UserAuthData{
		UserID:  userID,
		UserKey: authKeyStr,
	})
	response := struct {
		Code string `json:"code"`
	}{
		Code: authKeyStr,
	}
	return response, err
}

// getRecoveryMetrics retrieves a count of active vs expired recovery codes.
func getRecoveryMetrics(ctx context.Context, db *database.Connection, expiryInMinutes int64) (interface{}, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	query := sq.Select().
		Column(sq.Expr("COUNT(CASE WHEN t.age_in_minutes < ? THEN 1 END) AS active", expiryInMinutes)).
		Column(sq.Expr("COUNT(CASE WHEN t.age_in_minutes >= ? THEN 1 END) AS expired", expiryInMinutes)).
		FromSelect(sq.Select().
			Column(sq.Expr("TIMESTAMPDIFF(minute, created_at, ?) AS age_in_minutes", time.Now())).
			From("auth_scheme_data").
			Where(sq.Eq{"auth_scheme": recoveryConsts.Code}), "t")

	var metrics struct {
		ExpiredCount int64 `db:"expired" json:"expiredCount"`
		ActiveCount  int64 `db:"active" json:"activeCount"`
	}
	err := db.Get(&metrics, query)

	return metrics, err
}
