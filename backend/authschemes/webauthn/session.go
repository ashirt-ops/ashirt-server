package webauthn

import (
	"encoding/gob"

	auth "github.com/duo-labs/webauthn/webauthn"
)

type preRegistrationSessionData struct {
	UserData            webauthnUser
	WebAuthNSessionData *auth.SessionData
}

func init() {
	gob.Register(&preRegistrationSessionData{})
}
