package authschemes

// This is copied from: https://github.com/okta/okta-jwt-verifier-golang
// Copyright Okta, Inc, 2015-2018

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}
