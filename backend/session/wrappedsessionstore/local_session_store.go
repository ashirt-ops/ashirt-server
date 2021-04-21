// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

/*
Gorilla Sessions backend for MySQL.

Copyright (c) 2013 Contributors. See the list of contributors in the CONTRIBUTORS file for details.

This software is licensed under a MIT style license available in the LICENSE file.
*/

// Copyright (c) 2013 Gregor Robinson.
// Copyright (c) 2013 Brian Jones.
// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wrappedsessionstore

import (
	"database/sql"
	"database/sql/driver"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// DeletableSessionStore is an extension of sessions.Store that also supports a Delete method
type DeletableSessionStore interface {
	sessions.Store
	Delete(r *http.Request, w http.ResponseWriter, s *sessions.Session) error
}

// MySQLStore acts as a session store based on MySQL. This particular store also allows for
// identification of the user associated with the session
type MySQLStore struct {
	db         *sql.DB
	stmtInsert *sql.Stmt
	stmtDelete *sql.Stmt
	stmtUpdate *sql.Stmt
	stmtSelect *sql.Stmt

	Codecs          []securecookie.Codec
	Options         *sessions.Options
	table           string
	sessionToUserID func(*sessions.Session) *int64
}

type sessionRow struct {
	id         string
	userID     *int64
	data       string
	createdAt  time.Time
	modifiedAt time.Time
	expiresAt  time.Time
}

func init() {
	gob.Register(time.Time{})
}

// NewMySQLStore connects to the named database, then executes NewMySQLStoreFromConnection
func NewMySQLStore(endpoint string, path string, maxAge int, keyPairs ...[]byte) (*MySQLStore, error) {
	db, err := sql.Open("mysql", endpoint)
	if err != nil {
		return nil, err
	}

	return NewMySQLStoreFromConnection(db, path, maxAge, keyPairs...)
}

// NewMySQLStoreFromConnection prepares a few statements for interacting with the database, then returns
// a ready-to-use MySQLStore
func NewMySQLStoreFromConnection(db *sql.DB, path string, maxAge int, keyPairs ...[]byte) (*MySQLStore, error) {
	tableName := "sessions"

	insQ := "INSERT INTO sessions" +
		" (id, user_id, session_data, created_at, modified_at, expires_at) VALUES (NULL, ?, ?, ?, ?, ?)"
	stmtInsert, stmtErr := db.Prepare(insQ)
	if stmtErr != nil {
		return nil, stmtErr
	}

	delQ := "DELETE FROM sessions WHERE id = ?"
	stmtDelete, stmtErr := db.Prepare(delQ)
	if stmtErr != nil {
		return nil, stmtErr
	}

	updQ := "UPDATE sessions SET user_id=?, session_data = ?, created_at = ?, expires_at = ? WHERE id = ?"
	stmtUpdate, stmtErr := db.Prepare(updQ)
	if stmtErr != nil {
		return nil, stmtErr
	}

	selQ := "SELECT id, user_id, session_data, created_at, modified_at, expires_at from sessions WHERE id = ?"
	stmtSelect, stmtErr := db.Prepare(selQ)
	if stmtErr != nil {
		return nil, stmtErr
	}

	// mysqlstore still passes the values to securecookie for encryption before inserting into the database.
	// because of this, for large sessions (like ones containing session tokens we get from okta), securecookie
	// can fail to encrypt because it requires its output to be smaller to the max cookie size of 4096.
	// Since the output of securecookie is not actually going to a cookie, but instead to mysql we can remove this
	// cap by calling `MaxLength` on the codec.
	codecs := securecookie.CodecsFromPairs(keyPairs...)
	for _, codec := range codecs {
		if c, ok := codec.(*securecookie.SecureCookie); ok {
			c.MaxLength(0)
		}
	}

	return &MySQLStore{
		db:         db,
		stmtInsert: stmtInsert,
		stmtDelete: stmtDelete,
		stmtUpdate: stmtUpdate,
		stmtSelect: stmtSelect,
		Codecs:     codecs,
		Options: &sessions.Options{
			Path:   path,
			MaxAge: maxAge,
		},
		table: tableName,
	}, nil
}

// Close cleans up resources used by the MySQLStore (namely: statements, database access)
func (m *MySQLStore) Close() {
	m.stmtSelect.Close()
	m.stmtUpdate.Close()
	m.stmtDelete.Close()
	m.stmtInsert.Close()
	m.db.Close()
}

func (m *MySQLStore) SetSessionToUserID(fn func(*sessions.Session) *int64) {
	m.sessionToUserID = fn
}

func (m *MySQLStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(m, name)
}

func (m *MySQLStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(m, name)
	session.Options = &sessions.Options{
		Path:     m.Options.Path,
		Domain:   m.Options.Domain,
		MaxAge:   m.Options.MaxAge,
		Secure:   m.Options.Secure,
		HttpOnly: m.Options.HttpOnly,
	}
	session.IsNew = true
	var err error
	if cook, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, cook.Value, &session.ID, m.Codecs...)
		if err == nil {
			err = m.load(session)
			if err == nil {
				session.IsNew = false
			} else {
				err = nil
			}
		}
	}
	return session, err
}

func (m *MySQLStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	var err error
	if session.ID == "" {
		if err = m.insert(session); err != nil {
			return err
		}
	} else if err = m.save(session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, m.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

func (m *MySQLStore) insert(session *sessions.Session) error {
	var createdOn time.Time
	var modifiedOn time.Time
	var expiresOn time.Time
	crOn := session.Values["created_at"]
	if crOn == nil {
		createdOn = time.Now()
	} else {
		createdOn = crOn.(time.Time)
	}
	modifiedOn = createdOn
	exOn := session.Values["expires_at"]
	if exOn == nil {
		expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
	} else {
		expiresOn = exOn.(time.Time)
	}
	delete(session.Values, "created_at")
	delete(session.Values, "expires_at")
	delete(session.Values, "modified_at")

	var userID *int64
	if m.sessionToUserID != nil {
		userID = m.sessionToUserID(session)
	}

	encoded, encErr := securecookie.EncodeMulti(session.Name(), session.Values, m.Codecs...)
	if encErr != nil {
		return encErr
	}
	res, insErr := m.stmtInsert.Exec(userID, encoded, createdOn, modifiedOn, expiresOn)
	if insErr != nil {
		return insErr
	}
	lastInserted, lInsErr := res.LastInsertId()
	if lInsErr != nil {
		return lInsErr
	}
	session.ID = fmt.Sprintf("%d", lastInserted)
	return nil
}

// Delete removes the provided session from the store (database). An error is returned if some database
// issue occurs while trying to remove the session
func (m *MySQLStore) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Set cookie to expire.
	options := *session.Options
	options.MaxAge = -1
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", &options))
	// Clear session values.
	for k := range session.Values {
		delete(session.Values, k)
	}

	_, delErr := m.stmtDelete.Exec(session.ID)
	if delErr != nil {
		return delErr
	}
	return nil
}

func (m *MySQLStore) save(session *sessions.Session) error {
	if session.IsNew == true {
		return m.insert(session)
	}
	var createdOn time.Time
	var expiresOn time.Time
	crOn := session.Values["created_at"]
	if crOn == nil {
		createdOn = time.Now()
	} else {
		createdOn = crOn.(time.Time)
	}

	exOn := session.Values["expires_at"]
	if exOn == nil {
		expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
	} else {
		expiresOn = exOn.(time.Time)
		if expiresOn.Sub(time.Now().Add(time.Second*time.Duration(session.Options.MaxAge))) < 0 {
			expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
		}
	}

	delete(session.Values, "created_at")
	delete(session.Values, "expires_at")
	delete(session.Values, "modified_at")

	var userID *int64
	if m.sessionToUserID != nil {
		userID = m.sessionToUserID(session)
	}

	encoded, encErr := securecookie.EncodeMulti(session.Name(), session.Values, m.Codecs...)
	if encErr != nil {
		return encErr
	}
	_, updErr := m.stmtUpdate.Exec(userID, encoded, createdOn, expiresOn, session.ID)
	if updErr != nil {
		return updErr
	}
	return nil
}

func (m *MySQLStore) load(session *sessions.Session) error {
	row := m.stmtSelect.QueryRow(session.ID)
	sess := sessionRow{}
	var timeCreated, timeModified, timeExpires driver.Value
	scanErr := row.Scan(&sess.id, &sess.userID, &sess.data, &timeCreated, &timeModified, &timeExpires)
	if scanErr != nil {
		return scanErr
	}

	sess.createdAt = timeCreated.(time.Time)
	sess.modifiedAt = timeModified.(time.Time)
	sess.expiresAt = timeExpires.(time.Time)

	if sess.expiresAt.Sub(time.Now()) < 0 {
		return errors.New("Session expired")
	}
	err := securecookie.DecodeMulti(session.Name(), sess.data, &session.Values, m.Codecs...)
	if err != nil {
		return err
	}
	session.Values["created_at"] = sess.createdAt
	session.Values["modified_at"] = sess.modifiedAt
	session.Values["expires_at"] = sess.expiresAt
	return nil
}
