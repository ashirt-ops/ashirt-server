// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package authschemes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
)

func TestCreateNewAuthForUserGeneric(t *testing.T) {
	db, _, bridge := initBridgeTest(t)
	userID := createDummyUser(t, bridge, "normalUser")

	err := authschemes.CreateNewAuthForUserGeneric(db, "someauth", "someauth-type", authschemes.UserAuthData{
		UserID:  userID,
		UserKey: "dummy-user-key",
	})

	require.NoError(t, err)
}
