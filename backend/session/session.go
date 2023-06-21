package session

import (
	"encoding/gob"
)

// TODO TN should this be changed back to session?
type SessionData struct {
	UserID         int64
	IsAdmin        bool
	AuthSchemeData interface{}
}

func init() {
	gob.Register(&SessionData{})
}
