// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package webauthn

import "time"

type ListKeysOutput struct {
	Keys []KeyEntry `json:"keys"`
}

type KeyEntry struct {
	CredentialName string    `json:"credentialName"`
	DateCreated    time.Time `json:"dateCreated"`
}
