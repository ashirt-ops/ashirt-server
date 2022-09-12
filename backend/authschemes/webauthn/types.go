package webauthn

import (
	auth "github.com/duo-labs/webauthn/webauthn"
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
	KeyName             string
	UserID              int64
	RegistrationType    RegistrationType
	ExistingCredentials []AShirtWebauthnCredential
}

type AShirtWebauthnCredential struct {
	auth.Credential
	KeyName string `json:"keyName"`
}

func unwrapCredential(cred AShirtWebauthnCredential) auth.Credential {
	return cred.Credential
}

func wrapCredential(cred auth.Credential, keyName string) AShirtWebauthnCredential {
	return AShirtWebauthnCredential{
		Credential: cred,
		KeyName:    keyName,
	}
}
