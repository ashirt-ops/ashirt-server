package webauthn

import (
	"time"

	auth "github.com/go-webauthn/webauthn/webauthn"
)

type RegistrationType int

const (
	// CreateOrLinkCredential reflects the usecase where
	CreateCredential RegistrationType = iota
	LinkCredential
	AddCredential
)

type WebAuthnRegistrationInfo struct {
	Email                 string
	Username              string
	FirstName             string
	LastName              string
	CredentialName        string
	UserID                int64
	RegistrationType      RegistrationType
	ExistingCredentials   []AShirtWebauthnCredential
	CredentialCreatedDate time.Time
}

type WebAuthnUpdateCredentialInfo struct {
	UserID            int64
	CredentialName    string
	NewCredentialName string
}

type AShirtWebauthnExtension struct {
	CredentialName        string    `json:"credentialName"`
	CredentialCreatedDate time.Time `json:"credentialCreatedDate"`
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

type CredentialFlags struct {
	BackupEligible bool `json:"backupEligible"`
	BackupState    bool `json:"backupState"`
}

func (cf *CredentialFlags) Validate() error {
	if cf.BackupEligible && !cf.BackupState {
		return errors.New("BackupState must be true if BackupEligible is true")
	}
	return nil
}
