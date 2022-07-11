package helpers_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	recoveryHelpers "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/helpers"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/database/seeding"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
)

func TestGenerateRecoveryCodeForUser(t *testing.T) {
	db := setupDb(t)

	targetUserID := seeding.UserRon.ID
	code, err := recoveryHelpers.GenerateRecoveryCodeForUser(db, targetUserID)
	require.NoError(t, err)

	var recoveryAuthEntry models.AuthSchemeData
	db.Get(&recoveryAuthEntry, sq.Select("*").From("auth_scheme_data").Where(sq.Eq{
		"auth_scheme": recoveryConsts.Code,
		"user_id":     targetUserID,
	}))
	require.Equal(t, code, recoveryAuthEntry.UserKey)
}

func setupDb(t *testing.T) *database.Connection {
	db := seeding.InitTestWithOptions(t, seeding.TestOptions{
		DatabasePath: helpers.Ptr("../../../migrations"),
		DatabaseName: helpers.Ptr("recovery-auth-helpers-test-db"),
	})
	seeding.ApplySeeding(t, seeding.HarryPotterSeedData, db)

	return db
}
