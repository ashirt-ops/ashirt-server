package helpers

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ashirt-ops/ashirt-server/internal/authschemes"
	recoveryConsts "github.com/ashirt-ops/ashirt-server/internal/authschemes/recoveryauth/constants"
	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/errors"
)

const recoveryKeyLength = 40

// GenerateRecoveryCodeForUser provides a mechanism for non-recovery services to generate and register
// a recovery authentication for the provided UserID.
//
// returns the generatedCode if successful, and an error, if one was encountered
func GenerateRecoveryCodeForUser(db *database.Connection, userID int64) (string, error) {
	authKey := make([]byte, recoveryKeyLength)
	if _, err := rand.Read(authKey); err != nil {
		return "", errors.WrapError("Unable to generate random recovery key", err)
	}
	authKeyStr := hex.EncodeToString(authKey)

	err := GenerateNewRecoveryRecord(db, authschemes.UserAuthData{
		UserID:   userID,
		Username: authKeyStr,
	})

	return authKeyStr, err
}

// GenerateNewRecoveryRecord is a shorthand for CreateNewAuthForUserGeneric with the recovery code
// constant provide
func GenerateNewRecoveryRecord(db *database.Connection, userAuthData authschemes.UserAuthData) error {
	return authschemes.CreateNewAuthForUserGeneric(db, recoveryConsts.Code, recoveryConsts.Code, userAuthData)
}
