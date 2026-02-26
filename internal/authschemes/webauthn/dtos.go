package webauthn

import "time"

type ListCredentialsOutput struct {
	Credentials []CredentialEntry `json:"credentials"`
}

type CredentialEntry struct {
	CredentialName string    `json:"credentialName"`
	DateCreated    time.Time `json:"dateCreated"`
	CredentialID   string    `json:"credentialId"`
}
