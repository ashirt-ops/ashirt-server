// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

// TODO: these tests rely on service-level testing data, which needs to be migrated out of services
// (as a practical matter, having to recreate this functionality everywhere is burdensome at best)
// However, all of these tests were run at the service level prior, and work. These should be re-enabled
// once the testing data can be reused.

package recoveryauth_test

import (
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/database/seeding"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/session"
)

func initBridge(t *testing.T) authschemes.AShirtAuthBridge {
	db := seeding.InitTestWithOptions(t, seeding.TestOptions{
		DatabasePath: helpers.StringPtr("../../migrations"),
		DatabaseName: helpers.StringPtr("recovery-auth-test-db"),
	})
	seeding.ApplySeeding(t, seeding.HarryPotterSeedData, db)
	sessionStore, err := session.NewStore(db, session.StoreOptions{SessionDuration: time.Hour, Key: []byte{}})
	require.NoError(t, err)
	return authschemes.MakeAuthBridge(db, sessionStore, "local", "local")
}

func TestDeleteExpiredRecoveryCodes(t *testing.T) {
	bridge := initBridge(t)
	db := bridge.GetDatabase()
	normalUser := seeding.UserRon
	adminUser := seeding.UserDumbledore
	ctx := seeding.SimpleFullContext(normalUser)

	// note: there's something odd going on between go/mysql, where sometimes 1 hour is not
	// less than/equal to 1 hour (or 60 minutes). This causes random test failures, so setting to 59
	// minutes to avoid that mess.
	expirationDuration := int64(59)

	validKeyName := "valid"
	// add some recovery codes
	createDummyRecoveryRecord(t, db, "expired", normalUser.ID, 1*time.Hour)
	createDummyRecoveryRecord(t, db, validKeyName, normalUser.ID, 0)
	createDummyRecoveryRecord(t, db, "admin-expired", adminUser.ID, 2*time.Hour)
	createDummyRecoveryRecord(t, db, "also-expired", normalUser.ID, 60*time.Minute)

	require.Equal(t, 4, len(getRecoveryRecords(t, db)))

	// verify non-admins have no access (and no records are removed)
	err := recoveryauth.DeleteExpiredRecoveryCodes(ctx, db, expirationDuration)
	require.Error(t, err)
	require.Equal(t, 4, len(getRecoveryRecords(t, db)))

	// verify admins have access + effect works
	ctx = seeding.SimpleFullContext(adminUser)
	err = recoveryauth.DeleteExpiredRecoveryCodes(ctx, db, expirationDuration)
	require.NoError(t, err)

	recoveryRecords := getRecoveryRecords(t, db)

	require.Equal(t, 1, len(recoveryRecords))
	require.Equal(t, recoveryRecords[0].UserKey, validKeyName)
}

func createDummyRecoveryRecord(t *testing.T, db *database.Connection, key string, userID int64, age time.Duration) {
	_, err := db.Insert("auth_scheme_data", map[string]interface{}{
		"auth_scheme": recoveryConsts.Code,
		"auth_type":   recoveryConsts.Code,
		"user_key":    key,
		"user_id":     userID,
		"created_at":  time.Now().Add(-1 * age), // add negative time to emulate subtraction
	})
	require.NoError(t, err)
}

func getRecoveryRecords(t *testing.T, db *database.Connection) []models.AuthSchemeData {
	var recoveryRecords []models.AuthSchemeData
	err := db.Select(&recoveryRecords, sq.Select("*").
		From("auth_scheme_data").
		Where(sq.Eq{"auth_scheme": recoveryConsts.Code}))
	require.NoError(t, err)
	return recoveryRecords
}
