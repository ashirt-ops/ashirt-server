package webauthn

import (
	"encoding/gob"

	auth "github.com/duo-labs/webauthn/webauthn"
)

type webAuthNSessionData struct {
	UserData            webauthnUser
	WebAuthNSessionData *auth.SessionData
}

func init() {
	gob.Register(&webAuthNSessionData{})
}

func makeWebauthNSessionData(user webauthnUser, data *auth.SessionData) *webAuthNSessionData {
	sessionData := webAuthNSessionData{
		UserData:            user,
		WebAuthNSessionData: data,
	}
	return &sessionData
}
