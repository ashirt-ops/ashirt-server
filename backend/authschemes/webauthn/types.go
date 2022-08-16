package webauthn

import (
	auth "github.com/duo-labs/webauthn/webauthn"
)

type WebAuthnRegistrationInfo struct {
	Email     string
	FirstName string
	LastName  string
	KeyName   string
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
