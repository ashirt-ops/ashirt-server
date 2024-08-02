package webauthn

import (
	"encoding/gob"
	"github.com/go-webauthn/webauthn/metadata"
	auth "github.com/go-webauthn/webauthn/webauthn"
)

type webAuthNSessionData struct {
	UserData            webauthnUser
	WebAuthNSessionData *auth.SessionData
	Metadata            *metadata.Metadata
}

func init() {
	gob.Register(&webAuthNSessionData{})
}

func makeWebauthNSessionData(user webauthnUser, data *auth.SessionData, meta *metadata.Metadata) *webAuthNSessionData {
	sessionData := webAuthNSessionData{
		UserData:            user,
		WebAuthNSessionData: data,
		Metadata:            meta,
	}
	return &sessionData
}

func makeDiscoverableWebauthNSessionData(data *auth.SessionData, meta *metadata.Metadata) *webAuthNSessionData {
	sessionData := webAuthNSessionData{
		WebAuthNSessionData: data,
		Metadata:            meta,
	}
	return &sessionData
}
