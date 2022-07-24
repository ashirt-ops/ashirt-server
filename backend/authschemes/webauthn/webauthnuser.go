package webauthn

import (
	"strings"

	auth "github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
)


type webauthnUser struct {
	UserID      []byte
	UserName    string
	IconURL     string
	Credentials []auth.Credential

	firstName   string
	lastName    string
	email       string
}

func makeWebAuthnUser(firstName, lastName, email string) webauthnUser {
	return webauthnUser{
		UserID:      []byte(uuid.New().String()),
		UserName:    email,
	}
}

func (u *webauthnUser) WebAuthnID() []byte {
	return u.UserID
}

func (u *webauthnUser) WebAuthnName() string {
	return u.UserName
}

func (u *webauthnUser) WebAuthnDisplayName() string {
	return strings.Join([]string{u.firstName, u.lastName}, " ")
}

func (u *webauthnUser) WebAuthnIcon() string {
	return u.IconURL
}

func (u *webauthnUser) WebAuthnCredentials() []auth.Credential {
	return u.Credentials
}

func (u *webauthnUser) FirstName() string {
	return u.firstName
}

func (u *webauthnUser) LastName() string {
	return u.lastName
}

func (u *webauthnUser) Email() string {
	return u.email
}
