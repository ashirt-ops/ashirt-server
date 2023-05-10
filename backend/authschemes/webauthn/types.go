package webauthn

import (
	"time"

	auth "github.com/go-webauthn/webauthn/webauthn"
)

type RegistrationType int

const (
	// CreateOrLinkKey reflects the usecase where
	CreateKey RegistrationType = iota
	LinkKey
	AddKey
)

type WebAuthnRegistrationInfo struct {
	Email               string
	Username            string
	FirstName           string
	LastName            string
	CredentialName      string
	UserID              int64
	RegistrationType    RegistrationType
	ExistingCredentials []AShirtWebauthnCredential
	KeyCreatedDate      time.Time
}

type AShirtWebauthnExtension struct {
	CredentialName string    `json:"credentialName"`
	KeyCreatedDate time.Time `json:"keyCreatedDate"`
}

type AShirtWebauthnCredential struct {
	auth.Credential
	AShirtWebauthnExtension
}

func unwrapCredential(cred AShirtWebauthnCredential) auth.Credential {
	return cred.Credential
}

func wrapCredential(cred auth.Credential, extra AShirtWebauthnExtension) AShirtWebauthnCredential {
	return AShirtWebauthnCredential{
		Credential:              cred,
		AShirtWebauthnExtension: extra,
	}
}
