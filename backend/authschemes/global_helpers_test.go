package authschemes_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	"github.com/stretchr/testify/require"
)

func TestCreateNewAuthForUserGeneric(t *testing.T) {
	db, _, bridge := initBridgeTest(t)
	userID := createDummyUser(t, bridge, "normalUser")

	err := authschemes.CreateNewAuthForUserGeneric(db, "someauth", "someauth-type", authschemes.UserAuthData{
		UserID:   userID,
		Username: "dummy-user-key",
	})

	require.NoError(t, err)
}
