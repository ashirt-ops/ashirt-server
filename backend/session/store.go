// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package session

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/theparanoids/ashirt-server/backend/config"
)

// MySQLStore represents the session store.
type MySQLStore struct {
	*sql.DB
	stopCleanup chan bool
}

type SessionRow struct {
	Data string `json:"session_data"`
}

func GetSession(sessionManager *scs.SessionManager, r *http.Request) *Session {
	sessionData := sessionManager.Get(r.Context(), config.SessionStoreKey())

	if session, ok := sessionData.(*Session); ok {
		return session
	}

	return &Session{}
}

// New returns a new MySQLStore instance, with a background cleanup goroutine
// that runs every 5 minutes to remove expired session data.
func New(db *sql.DB) *MySQLStore {
	return NewWithCleanupInterval(db, 5*time.Minute)
}

// TODO TN do I want to use this???
// NewWithCleanupInterval returns a new MySQLStore instance. The cleanupInterval
// parameter controls how frequently expired session data is removed by the
// background cleanup goroutine. Setting it to 0 prevents the cleanup goroutine
// from running (i.e. expired sessions will not be removed).
func NewWithCleanupInterval(db *sql.DB, cleanupInterval time.Duration) *MySQLStore {
	m := &MySQLStore{
		DB: db,
	}

	if cleanupInterval > 0 {
		go m.startCleanup(cleanupInterval)
	}

	return m
}

// Find returns the data for a given session id from the MySQLStore instance.
// If the session id is not found or is expired, the returned exists flag will
// be set to false.

func (m *MySQLStore) Find(id string) ([]byte, bool, error) {
	sess := SessionRow{}
	row := m.DB.QueryRow("SELECT session_data FROM sessions WHERE id = ? AND UTC_TIMESTAMP(6) < expires_at", id)
	err := row.Scan(&sess.Data)
	if err == sql.ErrNoRows {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	fmt.Printf("sess: %+v\n", sess)
	fmt.Println("data: ", sess.Data)
	byteData := []byte(sess.Data)
	fmt.Println("byteData: ", byteData)

	return byteData, true, nil
}

func Decode(b []byte) (time.Time, map[string]interface{}, error) {
	aux := &struct {
		Deadline time.Time
		Values   map[string]interface{}
	}{}

	r := bytes.NewReader(b)
	if err := gob.NewDecoder(r).Decode(&aux); err != nil {
		return time.Time{}, nil, err
	}

	return aux.Deadline, aux.Values, nil
}

// Commit adds a session id and data to the MySQLStore instance with the given
// expiry time. If the session id already exists, then the data and expiry
// time are updated.
// TODO TN - what is current expirity time? and how to change new library to use it?
func (m *MySQLStore) Commit(id string, b []byte, expiry time.Time) error {
	_, val, err := Decode(b)
	// TODO TN add error message
	// fmt.Printf("val[config.SessionStoreKey()]: %+v\n", val[config.SessionStoreKey()].(*Session).UserID)
	userID := val[config.SessionStoreKey()].(*Session).UserID
	_, err = m.DB.Exec("INSERT INTO sessions (id, user_id, session_data, expires_at) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE session_data = VALUES(session_data), expires_at = VALUES(expires_at)", id, userID, b, expiry.UTC())
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a session id and corresponding data from the MySQLStore
// instance.
func (m *MySQLStore) Delete(id string) error {
	_, err := m.DB.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

func (m *MySQLStore) startCleanup(interval time.Duration) {
	m.stopCleanup = make(chan bool)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := m.deleteExpired()
			if err != nil {
				log.Println(err)
			}
		case <-m.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

// StopCleanup terminates the background cleanup goroutine for the MySQLStore
// instance. It's rare to terminate this; generally MySQLStore instances and
// their cleanup goroutines are intended to be long-lived and run for the lifetime
// of your application.
//
// There may be occasions though when your use of the MySQLStore is transient.
// An example is creating a new MySQLStore instance in a test function. In this
// scenario, the cleanup goroutine (which will run forever) will prevent the
// MySQLStore object from being garbage collected even after the test function
// has finished. You can prevent this by manually calling StopCleanup.
func (m *MySQLStore) StopCleanup() {
	if m.stopCleanup != nil {
		m.stopCleanup <- true
	}
}

func (m *MySQLStore) deleteExpired() error {
	_, err := m.DB.Exec("DELETE FROM sessions WHERE expires_at < UTC_TIMESTAMP(6)")
	return err
}
