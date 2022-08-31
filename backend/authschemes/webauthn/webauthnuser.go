package webauthn

import (
	"encoding/binary"
	"strings"

	auth "github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/theparanoids/ashirt-server/backend/helpers"
)

type webauthnUser struct {
	UserID      []byte
	UserName    string
	IconURL     string
	Credentials []AShirtWebauthnCredential
	FirstName   string
	LastName    string
	Email       string
	KeyName     string
}

func makeNewWebAuthnUser(firstName, lastName, email, keyName string) webauthnUser {
	return webauthnUser{
		UserID:    []byte(uuid.New().String()),
		UserName:  email,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		KeyName:   keyName,
	}
}

func makeLinkingWebAuthnUser(userID int64, email, keyName string) webauthnUser {
	return webauthnUser{
		UserID:   i64ToByteSlice(userID),
		UserName: email,
		Email:    email,
		KeyName:  keyName,
	}
}

func makeAddKeyWebAuthnUser(userID int64, email, keyName string, creds []AShirtWebauthnCredential) webauthnUser {
	user := makeLinkingWebAuthnUser(userID, email, keyName)
	user.Credentials = creds
	return user
}

func makeWebAuthnUser(firstName, lastName, slug, email string, userID int64, creds []AShirtWebauthnCredential) webauthnUser {
	return webauthnUser{
		UserID:      i64ToByteSlice(userID),
		UserName:    slug,
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
	return u.UserID
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
