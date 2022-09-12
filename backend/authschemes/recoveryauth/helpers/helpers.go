package helpers

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/database"
)

const recoveryKeyLength = 40

// GenerateRecoveryCodeForUser provides a mechanism for non-recovery services to generate and register
// a recovery authentication for the provided UserID.
//
// returns the generatedCode if successful, and an error, if one was encountered
func GenerateRecoveryCodeForUser(db *database.Connection, userID int64) (string, error) {
	authKey := make([]byte, recoveryKeyLength)
	if _, err := rand.Read(authKey); err != nil {
		return "", backend.WrapError("Unable to generate random recovery key", err)
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
