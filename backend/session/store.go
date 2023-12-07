package session

import (
	"net/http"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/session/wrappedsessionstore"
	"github.com/gorilla/sessions"
)

const sessionDataKey = "session_data"

type Store struct {
	wrappedStore wrappedsessionstore.DeletableSessionStore
}

type StoreOptions struct {
	SessionDuration  time.Duration
	UseSecureCookies bool
	Key              []byte
}

func NewStore(db *database.Connection, opts StoreOptions) (*Store, error) {
	wrappedStore, err := wrappedsessionstore.NewMySQLStoreFromConnection(
		db.DB,
		"/",
		int(opts.SessionDuration.Seconds()),
		opts.Key,
	)
	if err != nil {
		return nil, err
	}
	wrappedStore.SetSessionToUserID(userIDFromSession)
	wrappedStore.Options.Path = "/"
	wrappedStore.Options.HttpOnly = true
	wrappedStore.Options.Secure = opts.UseSecureCookies

	return &Store{wrappedStore}, nil
}

func (store *Store) Read(r *http.Request) *Session {
	sess := store.readRaw(r)
	sessionData, ok := sess.Values[sessionDataKey].(*Session)
	if ok {
		return sessionData
	}
	return &Session{}
}

func (store *Store) Set(w http.ResponseWriter, r *http.Request, s *Session) error {
	sess := store.readRaw(r)
	sess.Values[sessionDataKey] = s
	return sess.Save(r, w)
}

func (store *Store) Delete(w http.ResponseWriter, r *http.Request) error {
	return store.wrappedStore.Delete(r, w, store.readRaw(r))
}

func (store *Store) readRaw(r *http.Request) *sessions.Session {
	sess, _ := store.wrappedStore.Get(r, "auth") // ignoring, because errors only fire if a bad cookie name is provided
	return sess
}

func userIDFromSession(sess *sessions.Session) *int64 {
	sessionData, ok := sess.Values[sessionDataKey].(*Session)
	if ok {
		val := sessionData.UserID
		// userIDs start at 1. In some cases, the UserID won't be set and will be 0
		// trying to add a row with userID 0 will result in a foreign key constraint error, so
		// returning nil ("no user id") instead
		if val != 0 {
			return &val
		}
		return nil
	}
	return nil
}
