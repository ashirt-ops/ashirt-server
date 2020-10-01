// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package authschemes_test

import (
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/session"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func TestCreateNewUser(t *testing.T) {
	db, _, bridge := initBridgeTest(t)

	newUser, err := bridge.CreateNewUser(authschemes.UserProfile{
		FirstName: "Alice",
		LastName:  "Defaultuser",
		Email:     "alice@example.com",
		Slug:      "slug",
	})
	require.NoError(t, err)

	var user models.User
	getUserQuery := sq.Select("*").From("users")

	err = db.Get(&user, getUserQuery.Where(sq.Eq{"id": newUser.UserID}))
	require.NoError(t, err)
	require.Equal(t, "Alice", user.FirstName)
	require.Equal(t, "Defaultuser", user.LastName)
	require.Equal(t, "slug", user.Slug)

	// Creating a user with a slug that already exists appends a random number to the slug
	newUser, err = bridge.CreateNewUser(authschemes.UserProfile{
		FirstName: "Bob",
		LastName:  "Snooper",
		Email:     "bob@example.com",
		Slug:      "slug",
	})
	require.NoError(t, err)

	err = db.Get(&user, getUserQuery.Where(sq.Eq{"id": newUser.UserID}))
	require.NoError(t, err)
	require.Equal(t, "Bob", user.FirstName)
	require.Equal(t, "Snooper", user.LastName)
	require.Regexp(t, "slug-\\d{1,6}", user.Slug)
}

type testSession struct{ Some string }

func TestLoginUser(t *testing.T) {
	_, sessionStore, bridge := initBridgeTest(t)

	gob.Register(&testSession{})

	userID := createDummyUser(t, bridge)

	browser := &testBrowser{}
	w, r := browser.newRequest()
	err := bridge.LoginUser(w, r, userID, &testSession{Some: "data"})
	require.NoError(t, err)

	_, r = browser.newRequest()
	session := sessionStore.Read(r)
	require.NoError(t, err)
	require.Equal(t, userID, session.UserID)
	require.Equal(t, "data", session.AuthSchemeData.(*testSession).Some)
}

func TestAddToSession(t *testing.T) {
	_, _, bridge := initBridgeTest(t)

	gob.Register(&testSession{})

	browser := &testBrowser{}
	w, r := browser.newRequest()
	bridge.SetAuthSchemeSession(w, r, &testSession{Some: "data"})

	_, r = browser.newRequest()
	require.Equal(t, &testSession{Some: "data"}, bridge.ReadAuthSchemeSession(r))
}

func TestDeleteSession(t *testing.T) {
	_, sessionStore, bridge := initBridgeTest(t)

	gob.Register(&testSession{})

	userID := createDummyUser(t, bridge)

	browser := &testBrowser{}
	w, r := browser.newRequest()
	err := bridge.LoginUser(w, r, userID, &testSession{Some: "data"})
	require.NoError(t, err)

	w, r = browser.newRequest()
	bridge.DeleteSession(w, r)

	_, r = browser.newRequest()
	session := sessionStore.Read(r)
	require.Equal(t, int64(0), session.UserID)
}

func TestUserAuthCreationAndLookup(t *testing.T) {
	_, _, bridge := initBridgeTest(t)

	userID := createDummyUser(t, bridge)
	err := bridge.CreateNewAuthForUser(authschemes.UserAuthData{
		UserID:  userID,
		UserKey: "dummy-user-key",
	})
	require.NoError(t, err)

	t.Run("Test FindUserAuth", func(t *testing.T) {
		auth, err := bridge.FindUserAuth("dummy-user-key")
		require.NoError(t, err)
		require.Equal(t, userID, auth.UserID)
		require.Equal(t, "dummy-user-key", auth.UserKey)
	})

	t.Run("Test FindUserAuthsByUserSlug", func(t *testing.T) {
		auths, err := bridge.FindUserAuthsByUserSlug("dummy-user-slug")
		require.NoError(t, err)
		require.Len(t, auths, 1)
		require.Equal(t, userID, auths[0].UserID)
		require.Equal(t, "dummy-user-key", auths[0].UserKey)
	})

	t.Run("Test UpdateAuthForUser", func(t *testing.T) {
		authData := authschemes.UserAuthData{
			UserKey:            "dummy-user-key",
			EncryptedPassword:  []byte("encrypted-password"),
			NeedsPasswordReset: true,
		}
		err := bridge.UpdateAuthForUser(authData)
		require.NoError(t, err)

		auth, err := bridge.FindUserAuth("dummy-user-key")
		require.NoError(t, err)
		require.Equal(t, []byte("encrypted-password"), auth.EncryptedPassword)
		require.Equal(t, true, auth.NeedsPasswordReset)
	})
}

func initBridgeTest(t *testing.T) (*database.Connection, *session.Store, authschemes.AShirtAuthBridge) {
	db := database.NewTestConnection(t, "authschemes-test-db")
	sessionStore, err := session.NewStore(db, session.StoreOptions{SessionDuration: time.Hour, Key: []byte{}})
	require.NoError(t, err)
	return db, sessionStore, authschemes.MakeAuthBridge(db, sessionStore, "test")
}

func createDummyUser(t *testing.T, bridge authschemes.AShirtAuthBridge) int64 {
	newUser, err := bridge.CreateNewUser(authschemes.UserProfile{
		FirstName: "Dummy",
		LastName:  "User",
		Email:     "email@example.com",
		Slug:      "dummy-user-slug",
	})
	require.NoError(t, err)
	return newUser.UserID
}

// testBrowser generates test requests/responsewriters and saves cookies for all future requests
type testBrowser struct {
	lastResponseRecorder *httptest.ResponseRecorder
	cookies              []*http.Cookie
}

func (b *testBrowser) newRequest() (http.ResponseWriter, *http.Request) {
	// Save cookies from last recorded response
	if b.lastResponseRecorder != nil {
		cookiesToAdd := b.lastResponseRecorder.Result().Cookies()
		for _, cookie := range cookiesToAdd {
			b.cookies = append(b.cookies, cookie)
		}
	}

	r := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()

	// Add all saved cookies to the request
	for _, cookie := range b.cookies {
		r.AddCookie(cookie)
	}

	b.lastResponseRecorder = responseRecorder

	return responseRecorder, r
}
