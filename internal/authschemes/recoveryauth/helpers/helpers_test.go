package helpers_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	recoveryConsts "github.com/ashirt-ops/ashirt-server/internal/authschemes/recoveryauth/constants"
	recoveryHelpers "github.com/ashirt-ops/ashirt-server/internal/authschemes/recoveryauth/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/database/seeding"
	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/stretchr/testify/require"
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
	require.Equal(t, code, recoveryAuthEntry.Username)
}

func setupDb(t *testing.T) *database.Connection {
	db := seeding.InitTestWithOptions(t, seeding.TestOptions{
		DatabasePath: helpers.Ptr("../../../migrations"),
		DatabaseName: helpers.Ptr("recovery-auth-helpers-test-db"),
	})
	seeding.ApplySeeding(t, seeding.HarryPotterSeedData, db)

	return db
}
