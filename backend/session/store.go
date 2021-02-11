// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package session

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/session/wrappedsessionstore"
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
		return &sessionData.UserID
	}
	return nil
}
