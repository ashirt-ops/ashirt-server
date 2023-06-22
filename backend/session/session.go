package session

import (
	"encoding/gob"
)

type Session struct {
	UserID         int64
	IsAdmin        bool
	AuthSchemeData interface{}
}

func init() {
	gob.Register(&Session{})
}
