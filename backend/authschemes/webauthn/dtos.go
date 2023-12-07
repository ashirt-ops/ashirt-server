// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

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
