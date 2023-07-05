package session

import (
	"encoding/gob"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/theparanoids/ashirt-server/backend/config"
)

type Session struct {
	UserID         int64
	IsAdmin        bool
	AuthSchemeData interface{}
}

func init() {
	gob.Register(&Session{})
}

func GetSession(sessionManager *scs.SessionManager, r *http.Request) *Session {
	sessionData := sessionManager.Get(r.Context(), config.SessionStoreKey())

	if session, ok := sessionData.(*Session); ok {
		return session
	}

	return &Session{}
}
