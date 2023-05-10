package webauthn

import (
	"encoding/binary"
	"strings"
	"time"

	auth "github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/theparanoids/ashirt-server/backend/helpers"
)

type webauthnUser struct {
	UserID         []byte
	AuthnID        []byte
	UserName       string
	IconURL        string
	Credentials    []AShirtWebauthnCredential
	FirstName      string
	LastName       string
	Email          string
	CredentialName string
	KeyCreatedDate time.Time
}

func makeNewWebAuthnUser(firstName, lastName, email, username, credentialName string) webauthnUser {
	return webauthnUser{
		AuthnID:        []byte(uuid.New().String()),
		UserName:       username,
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		CredentialName: credentialName,
		KeyCreatedDate: time.Now(),
	}
}

func makeLinkingWebAuthnUser(userID int64, username, credentialName string) webauthnUser {
	return webauthnUser{
		UserID:         i64ToByteSlice(userID),
		AuthnID:        []byte(uuid.New().String()),
		UserName:       username,
		CredentialName: credentialName,
		KeyCreatedDate: time.Now(),
	}
}

func makeAddKeyWebAuthnUser(userID int64, username, credentialName string, creds []AShirtWebauthnCredential) webauthnUser {
	user := makeLinkingWebAuthnUser(userID, username, credentialName)
	user.Credentials = creds
	return user
}

func makeWebAuthnUser(firstName, lastName, username, email string, UserID int64, authnID []byte, creds []AShirtWebauthnCredential) webauthnUser {
	return webauthnUser{
		AuthnID:     authnID,
		UserID:      i64ToByteSlice(UserID),
		UserName:    username,
		Credentials: creds,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
	}
}

func i64ToByteSlice(i int64) []byte {
	uInt := uint64(i)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uInt)
	return b
}

func byteSliceToI64(b []byte) int64 {
	uInt := binary.LittleEndian.Uint64(b)
	return int64(uInt)
}

func (u *webauthnUser) WebAuthnID() []byte {
	return u.AuthnID
}

func (u *webauthnUser) WebAuthnName() string {
	return u.UserName
}

func (u *webauthnUser) WebAuthnDisplayName() string {
	return strings.Join([]string{u.FirstName, u.LastName}, " ")
}

func (u *webauthnUser) WebAuthnIcon() string {
	return u.IconURL
}

func (u *webauthnUser) WebAuthnCredentials() []auth.Credential {
	return helpers.Map(u.Credentials, unwrapCredential)
}

func (u *webauthnUser) UserIDAsI64() int64 {
	return byteSliceToI64(u.UserID)
}
