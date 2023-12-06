package authschemes

import (
	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
)

// CreateNewAuthForUserGeneric provides a mechanism for non-auth providers to generate new authentications
// on behalf of auth providers. This is only intended for recovery.
//
// Proper usage:  authschemes.CreateNewAuthForUser(db, recoveryauth.constants.Code, authschemes.UserAuthData{})
// note: you will need to provide your own database instance
func CreateNewAuthForUserGeneric(db *database.Connection, authSchemeName, authSchemeType string, data UserAuthData) error {
	_, err := db.Insert("auth_scheme_data", map[string]interface{}{
		"auth_scheme":         authSchemeName,
		"auth_type":           authSchemeType,
		"username":            data.Username,
		"user_id":             data.UserID,
		"authn_id":            string(data.AuthnID),
		"encrypted_password":  data.EncryptedPassword,
		"totp_secret":         data.TOTPSecret,
		"must_reset_password": data.NeedsPasswordReset,
		"json_data":           data.JSONData,
	})
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return backend.BadInputErr(err, "An account for this user already exists")
		}
		return backend.WrapError("Unable to generate auth scheme for user", backend.DatabaseErr(err))
	}
	return nil
}
