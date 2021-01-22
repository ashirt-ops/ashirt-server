// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package recoveryauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

// DeleteExpiredRecoveryCodes removes recovery codes from the database that are older than X minutes old
// Duration is determined by looking at the environment configuration
func DeleteExpiredRecoveryCodes(ctx context.Context, db *database.Connection, expiryInMinutes int64) error {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Insufficient access to remove recovery code", backend.UnauthorizedWriteErr(err))
	}

	err := db.Delete(sq.Delete("auth_scheme_data").
		Where(sq.Eq{"auth_scheme": recoveryConsts.Code}).
		Where("TIMESTAMPDIFF(minute, created_at, ?) >= ?", time.Now(), expiryInMinutes)) // ensure timestamps are sufficently old
	if err != nil {
		return backend.WrapError("Unable to remove recovery code", backend.DatabaseErr(err))
	}

	return nil
}

// generateRecoveryCodeForUser creates a new, cryptographically random, hex-encoded,
// recovery key. This key is then attached to a user as a authorization method.
func generateRecoveryCodeForUser(ctx context.Context, bridge authschemes.AShirtAuthBridge, userSlug string) (interface{}, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.WrapError("Non-admin tried to generate recovery code", backend.UnauthorizedWriteErr(err))
	}

	userID, err := bridge.GetUserIDFromSlug(userSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to get UserID from slug", err)
	}

	authKeyStr, err := GenerateRecoveryCodeForUser(bridge.GetDatabase(), userID)
	if err != nil {
		return nil, backend.WrapError("Could not generate recovery code for user", err)
	}

	response := struct {
		Code string `json:"code"`
	}{Code: authKeyStr}

	if err != nil {
		return response, backend.WrapError("Unable to create recovery auth for user", err)
	}
	return response, nil
}

// getRecoveryMetrics retrieves a count of active vs expired recovery codes.
func getRecoveryMetrics(ctx context.Context, db *database.Connection, expiryInMinutes int64) (interface{}, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.WrapError("Non-admin tried to get recovery metrics", backend.UnauthorizedReadErr(err))
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
	if err != nil {
		return metrics, backend.WrapError("Failed to query recovery metrics", err)
	}

	return metrics, nil
}

func generateRecoveryEmail(ctx context.Context, bridge authschemes.AShirtAuthBridge, userEmail string) error {
	userAccounts, err := bridge.FindUserAuthsByUserEmail(userEmail)
	if err != nil || len(userAccounts) == 0 {
		if len(userAccounts) == 0 {
			err = backend.DatabaseErr(errors.New("No matching user accounts"))
		}
		return backend.WrapError("Unable to get user accounts from email", err)
	}

	targetAccount := userAccounts[0]
	err = bridge.AddScheduledEmail(userEmail, &targetAccount, "self-service-recovery-email")
	if err != nil {
		return backend.WrapError("Unable to generate recovery email", err)
	}

	return nil
}

// GenerateRecoveryCodeForUser provides a mechanism for non-recovery services to generate and register
// a recovery authentication for the provided UserID.
//
// returns the generatedCode if successful, and an error, if one was encountered
func GenerateRecoveryCodeForUser(db *database.Connection, userID int64) (string, error) {
	authKey := make([]byte, recoveryKeyLength)
	if _, err := rand.Read(authKey); err != nil {
		return "", backend.WrapError("Unable to generate random recovery key", err)
	}
	authKeyStr := hex.EncodeToString(authKey)

	err := GenerateNewRecoveryRecord(db, authschemes.UserAuthData{
		UserID:  userID,
		UserKey: authKeyStr,
	})

	return authKeyStr, err
}

// GenerateNewRecoveryRecord is a shorthand for CreateNewAuthForUserGeneric with the recovery code
// constant provide
func GenerateNewRecoveryRecord(db *database.Connection, userAuthData authschemes.UserAuthData) error {
	return authschemes.CreateNewAuthForUserGeneric(db, recoveryConsts.Code, userAuthData)
}
