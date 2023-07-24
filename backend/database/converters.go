// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"github.com/ashirt-ops/ashirt-server/backend/models"
)

// EvidenceToID is a small helper to grab the ID from a models.Evidence. Useful when paired with
// helpers.Map
func EvidenceToID(e models.Evidence) int64 {
	return e.ID
}
