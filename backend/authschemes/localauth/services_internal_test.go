package localauth

import (
	"context"
	"testing"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	"github.com/ashirt-ops/ashirt-server/backend/database/seeding"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/session"
	"github.com/stretchr/testify/require"
)

func initBridge(t *testing.T) authschemes.AShirtAuthBridge {
	db := seeding.InitTestWithOptions(t, seeding.TestOptions{
		DatabasePath: helpers.Ptr("../../migrations"),
		DatabaseName: helpers.Ptr("local-auth-test-db"),
	})
	seeding.ApplySeeding(t, seeding.HarryPotterSeedData, db)
	sessionStore, err := session.NewStore(db, session.StoreOptions{SessionDuration: time.Hour, Key: []byte{}})
	require.NoError(t, err)
	return authschemes.MakeAuthBridge(db, sessionStore, "local", "local")
}

func TestReadUserTotpStatus(t *testing.T) {
	bridge := initBridge(t)

	deviousUser := seeding.UserDraco
	adminUser := seeding.UserDumbledore
	targetUser := seeding.UserRon

	var ctx context.Context
	var hasTotp bool
	var err error

	// verify target user ("self") does not have totp
	ctx = seeding.SimpleFullContext(targetUser)
	hasTotp, err = readUserTotpStatus(ctx, bridge, targetUser.Slug) // with slug
	require.NoError(t, err)
	require.False(t, hasTotp)

	ctx = seeding.SimpleFullContext(targetUser)
	hasTotp, err = readUserTotpStatus(ctx, bridge, "") // without slug
	require.NoError(t, err)
	require.False(t, hasTotp)

	// verify that devious user cannot check totp status for target user
	ctx = seeding.SimpleFullContext(deviousUser)
	hasTotp, err = readUserTotpStatus(ctx, bridge, targetUser.Slug)
	require.Error(t, err)

	// verify admin _can_ check totp status for target user
	ctx = seeding.SimpleFullContext(adminUser)
	hasTotp, err = readUserTotpStatus(ctx, bridge, targetUser.Slug)
	require.NoError(t, err)
	require.False(t, hasTotp)

	// give target user totp -- note: this will remove existing encrypted_password and must_reset_password values
	err = bridge.UpdateAuthForUser(authschemes.UserAuthData{
		TOTPSecret: helpers.Ptr("abc123"),
		Username:   targetUser.FirstName,
	})
	require.NoError(t, err)

	// verify that target now has totp
	ctx = seeding.SimpleFullContext(targetUser)
	hasTotp, err = readUserTotpStatus(ctx, bridge, "")
	require.NoError(t, err)
	require.True(t, hasTotp)

	// verify admin sees totp change status for target user
	ctx = seeding.SimpleFullContext(adminUser)
	hasTotp, err = readUserTotpStatus(ctx, bridge, targetUser.Slug)
	require.NoError(t, err)
	require.True(t, hasTotp)
}

func TestDeleteUserTotp(t *testing.T) {
	bridge := initBridge(t)

	deviousUser := seeding.UserDraco
	adminUser := seeding.UserDumbledore
	targetUser := seeding.UserRon

	var ctx context.Context
	var err error

	// give target user TOTP
	err = bridge.UpdateAuthForUser(authschemes.UserAuthData{
		TOTPSecret: helpers.Ptr("abc123"),
		Username:   targetUser.FirstName,
	})
	require.NoError(t, err)

	// verify that a devious user cannot remove totp for a target user
	ctx = seeding.SimpleFullContext(deviousUser)
	err = deleteUserTotp(ctx, bridge, targetUser.Slug)
	require.Error(t, err)

	// verify target user can remove their own totp
	ctx = seeding.SimpleFullContext(targetUser)
	err = deleteUserTotp(ctx, bridge, "")
	require.NoError(t, err)

	// verify admin cannot "remove" target user totp (for user that does not have totp)
	ctx = seeding.SimpleFullContext(adminUser)
	err = deleteUserTotp(ctx, bridge, targetUser.Slug)
	require.Error(t, err)

	// re-give user totp
	err = bridge.UpdateAuthForUser(authschemes.UserAuthData{
		TOTPSecret: helpers.Ptr("abc123"),
		Username:   targetUser.FirstName,
	})
	require.NoError(t, err)

	// verify admin can remove target user totp
	ctx = seeding.SimpleFullContext(adminUser)
	err = deleteUserTotp(ctx, bridge, targetUser.Slug)
	require.NoError(t, err)
}
