// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package authschemes

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/services"
	"github.com/theparanoids/ashirt-server/backend/session"

	sq "github.com/Masterminds/squirrel"
)

// AShirtAuthBridge provides a set of functionality that bridges the identity resolution
// (the AuthScheme) and persistent user/session management
type AShirtAuthBridge struct {
	db             *database.Connection
	sessionStore   *session.Store
	authSchemeName string
	authSchemeType string
}

// MakeAuthBridge constructs returns a set of functions to interact with the underlying AShirt
// authentication scheme
func MakeAuthBridge(db *database.Connection, sessionStore *session.Store, authSchemeName, authSchemeType string) AShirtAuthBridge {
	return AShirtAuthBridge{
		db:             db,
		sessionStore:   sessionStore,
		authSchemeName: authSchemeName,
		authSchemeType: authSchemeType,
	}
}

// CreateNewUser allows new users to be registered into the system, if they do not already exist.
// Note that slug must be unique
func (ah AShirtAuthBridge) CreateNewUser(profile UserProfile) (*dtos.CreateUserOutput, error) {
	return services.CreateUser(ah.db, profile.ToCreateUserInput())
}

// SetAuthSchemeSession sets authscheme specific session data to the current user session. Session data should
// be a struct and registered with `gob.Register` in an init function of the authscheme
func (ah AShirtAuthBridge) SetAuthSchemeSession(w http.ResponseWriter, r *http.Request, data interface{}) error {
	s := ah.sessionStore.Read(r)
	s.AuthSchemeData = data
	return ah.sessionStore.Set(w, r, s)
}

// ReadAuthSchemeSession retrieves previously saved session data set by SetAuthSchemeSession
func (ah AShirtAuthBridge) ReadAuthSchemeSession(r *http.Request) interface{} {
	return ah.sessionStore.Read(r).AuthSchemeData
}

// LoginUser denotes that a user shall be logged in.
// In addition to the required userID, a user can also provide custom authscheme specific session data
func (ah AShirtAuthBridge) LoginUser(w http.ResponseWriter, r *http.Request, userID int64, authSchemeSessionData interface{}) error {
	if !(ah.isAccountEnabled(r, userID)) {
		ah.DeleteSession(w, r)
		return backend.WrapError("Unable to login user", backend.AccountDisabled())
	}

	ah.updateLastLogin(r, userID)

	return ah.sessionStore.Set(w, r, &session.Session{
		UserID:         userID,
		IsAdmin:        ah.isAdmin(r, userID),
		AuthSchemeData: authSchemeSessionData,
	})
}

func (ah AShirtAuthBridge) updateLastLogin(r *http.Request, userID int64) {
	err := ah.db.Update(sq.Update("auth_scheme_data").Set("last_login", time.Now()).Where(sq.Eq{"user_id": userID, "auth_scheme": ah.authSchemeName}))
	if err != nil {
		logging.Log(r.Context(), "msg", "Unable to update last_login", "userID", userID, "error", err)
	}
}

// IsAccountEnabled checks if the provided userid has an enabled account (specifically, it does not
// have the disabled flag set)
// returns (false, err) if the account cannot be found or another database error occurred.
func (ah AShirtAuthBridge) IsAccountEnabled(userID int64) (bool, error) {
	var flag bool
	err := ah.db.Get(&flag, sq.Select("disabled").From("users").Where(sq.Eq{"id": userID}))
	if err != nil {
		return false, err
	}

	return !flag, nil
}

func (ah AShirtAuthBridge) isAccountEnabled(r *http.Request, userID int64) bool {
	enabled, err := ah.IsAccountEnabled(userID)
	if err != nil {
		logging.Log(r.Context(), "msg", "Unable to check user's disabled flag", "userID", userID, "error", err)
		return false
	}
	return enabled
}

func (ah AShirtAuthBridge) isAdmin(r *http.Request, userID int64) bool {
	var isAdmin bool
	err := ah.db.Get(&isAdmin, sq.Select("admin").From("users").Where(sq.Eq{"id": userID}))
	if err != nil {
		logging.Log(r.Context(), "msg", "Unable to check user's admin flag", "userID", userID, "error", err)
		return false
	}
	return isAdmin
}

// GetUserIDFromSlug retrieves a userid from the provided user slug.
func (ah AShirtAuthBridge) GetUserIDFromSlug(userSlug string) (int64, error) {
	return ah.db.RetrieveUserIDBySlug(userSlug)
}

func (ah AShirtAuthBridge) GetUserFromID(userID int64) (models.User, error) {
	return ah.db.RetrieveUserByID(userID)
}

// DeleteSession removes a user's session. Useful in situtations where authentication fails,
// and we want to treat the user as not-logged-in
func (ah AShirtAuthBridge) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	return ah.sessionStore.Delete(w, r)
}

// UserAuthData is a small structure capturing data relevant to a user for authentication purposes
type UserAuthData struct {
	UserID             int64   `db:"user_id"`
	UserKey            string  `db:"user_key"`
	EncryptedPassword  []byte  `db:"encrypted_password"`
	NeedsPasswordReset bool    `db:"must_reset_password"`
	TOTPSecret         *string `db:"totp_secret"`
	JSONData           *string `db:"json_data"`
}

// FindUserAuth retrieves the row (codified by UserAuthData) corresponding to the provided userKey(e.g. username, email, etc) and the
// auth scheme name provided from the caller.
//
// Returns a fully populated UserAuthData object, or an error if no such row exists
func (ah AShirtAuthBridge) FindUserAuth(userKey string) (UserAuthData, error) {
	var authData UserAuthData

	err := ah.db.Get(&authData, sq.Select(
		"user_id", "user_key", "encrypted_password",
		"must_reset_password", "totp_secret", "json_data").
		From("auth_scheme_data").
		Where(sq.Eq{
			"user_key":    userKey,
			"auth_scheme": ah.authSchemeName,
		}))
	if err != nil {
		return UserAuthData{}, backend.WrapError("Cannot find user authentication", backend.DatabaseErr(err))
	}
	return authData, nil
}

// FindUserAuthByContext acts as a proxy to calling FindUserByUserID with the userID extracted from the provided context
//  see FindUserAuthByUserID
func (ah AShirtAuthBridge) FindUserAuthByContext(ctx context.Context) (UserAuthData, error) {
	return ah.FindUserAuthByUserID(middleware.UserID(ctx))
}

// FindUserAuthByUserID retrieves the row (codified by UserAuthData) corresponding to the provided userID
//
// Returns a fully populated UserAuthData object, or nil if no such row exists
func (ah AShirtAuthBridge) FindUserAuthByUserID(userID int64) (UserAuthData, error) {
	var authData UserAuthData

	err := ah.db.Get(&authData, sq.Select("user_id", "user_key", "encrypted_password", "must_reset_password", "totp_secret", "json_data").
		From("auth_scheme_data").
		Where(sq.Eq{
			"user_id":     userID,
			"auth_scheme": ah.authSchemeName,
		}))
	if err != nil {
		return UserAuthData{}, backend.DatabaseErr(err)
	}
	return authData, nil
}

func (ah AShirtAuthBridge) FindAuthsForUsername(userKey string) (UserAuthData, error) {
	var authData UserAuthData

	err := ah.db.Get(&authData, sq.Select("user_id", "user_key", "encrypted_password", "must_reset_password", "totp_secret", "json_data").
		From("auth_scheme_data").
		Where(sq.Eq{
			"user_key":    userKey,
			"auth_scheme": ah.authSchemeName,
		}))
	if err != nil {
		return UserAuthData{}, backend.DatabaseErr(err)
	}
	return authData, nil
}

func (ah AShirtAuthBridge) findUserAuthsByUserEmail(email string, includeDeleted bool) ([]UserAuthData, error) {
	var authData []UserAuthData

	whereClause := sq.Eq{
		"users.email": email,
	}
	if !includeDeleted {
		whereClause["users.deleted_at"] = nil
	}

	err := ah.db.Select(&authData, sq.Select("user_id", "user_key", "encrypted_password", "must_reset_password", "totp_secret", "json_data").
		From("auth_scheme_data").
		LeftJoin("users ON users.id = auth_scheme_data.user_id").
		Where(whereClause))
	if err != nil {
		return []UserAuthData{}, backend.DatabaseErr(err)
	}
	return authData, nil
}

// FindUserByEmail retrieves the user record associated with the specified email address. Returns
// the found record, or an error if no such record exists
func (ah AShirtAuthBridge) FindUserByEmail(email string, includeDeleted bool) (models.User, error) {
	var userRecord models.User

	whereClause := sq.Eq{
		"email": email,
	}
	if !includeDeleted {
		whereClause["deleted_at"] = nil
	}
	err := ah.db.Get(&userRecord,
		sq.Select("*").From("users").Where(whereClause))
	if err != nil {
		return models.User{}, backend.DatabaseErr(err)
	}
	return userRecord, nil
}

// CheckIfUserEmailTaken attempts to find an account with the provided email. If found, returns true,
// otherwise, returns false. Note that if the user found is the one with the ID matching allowUserID, then false is returned.
func (ah AShirtAuthBridge) CheckIfUserEmailTaken(email string, allowUserID int64, includeDeleted bool) (bool, error) {
	user, err := ah.FindUserByEmail(email, includeDeleted)
	if err == nil {
		if user.ID == allowUserID { // check if the user requesting the email account already "owns" it.
			return false, nil
		}
		return true, nil
	}

	trueErr, ok := err.(*backend.HTTPError)
	if ok && trueErr.WrappedError == sql.ErrNoRows {
		return false, nil
	}
	return false, err
}

// FindUserAuthsByUserEmail retrieves the rows (codified by UserAuthData) corresponding to the provided userEmail for NON-DELETED accounts.
// Note that a user may have multiple authentications based on a single email, so each of these records are returned.
//
// See FindUserAuthsByUserEmailIncludeDeleted to retreive all users irrespective of if they have been deleted
// Returns a fully populated UserAuthData object, or nil if no such row exists
func (ah AShirtAuthBridge) FindUserAuthsByUserEmail(email string) ([]UserAuthData, error) {
	return ah.findUserAuthsByUserEmail(email, false)
}

// FindUserAuthsByUserEmailIncludeDeleted retrieves the rows (codified by UserAuthData) corresponding to the provided userEmail for ALL accounts.
// Note that a user may have multiple authentications based on a single email, so each of these records are returned.
//
// Returns a fully populated UserAuthData object, or nil if no such row exists
func (ah AShirtAuthBridge) FindUserAuthsByUserEmailIncludeDeleted(email string) ([]UserAuthData, error) {
	return ah.findUserAuthsByUserEmail(email, true)
}

// FindUserAuthsByUserSlug retrieves the row (codified by UserAuthData) corresponding to the provided user slug and the
// auth scheme name provided from the caller.
//
// Returns a fully populated UserAuthData object, or nil if no such row exists
func (ah AShirtAuthBridge) FindUserAuthsByUserSlug(slug string) ([]UserAuthData, error) {
	var authData []UserAuthData

	err := ah.db.Select(&authData, sq.Select("user_id", "user_key", "encrypted_password", "must_reset_password", "totp_secret", "json_data").
		From("auth_scheme_data").
		LeftJoin("users ON users.id = auth_scheme_data.user_id").
		Where(sq.Eq{
			"users.slug":  slug,
			"auth_scheme": ah.authSchemeName,
		}))
	if err != nil {
		return []UserAuthData{}, backend.WrapError("Unable to fetch user authentications", backend.DatabaseErr(err))
	}
	return authData, nil
}

// CreateNewAuthForUser adds a new entry to the auth_scheme_data table for the given UserAuthData.
//
// Returns nil if no error was occurred, BadInputErr if the user account already exists, or DatabaseErr
// if any other issue occurs
func (ah AShirtAuthBridge) CreateNewAuthForUser(data UserAuthData) error {
	return CreateNewAuthForUserGeneric(ah.db, ah.authSchemeName, ah.authSchemeType, data)
}

// UpdateAuthForUser updates a user's authentication password, and can flag whether the user needs to
// change their password on the next login.
func (ah AShirtAuthBridge) UpdateAuthForUser(data UserAuthData) error {
	ub := sq.Update("auth_scheme_data").
		SetMap(map[string]interface{}{
			"encrypted_password":  data.EncryptedPassword,
			"must_reset_password": data.NeedsPasswordReset,
			"totp_secret":         data.TOTPSecret,
			"json_data":           data.JSONData,
		}).
		Where(sq.Eq{"user_key": data.UserKey, "auth_scheme": ah.authSchemeName})
	err := ah.db.Update(ub)
	if err != nil {
		return backend.WrapError("Unable to update user authentication", backend.DatabaseErr(err))
	}
	return nil
}

// OneTimeVerification looks for a matching record in the auth_scheme_data table with the following conditions:
// user_key matches && created_at less than <expirationInMinutes> minutes
// If this record exists, then the record is deleted. If there is no error _either_ for the lookup
// OR the deletion, then (userID for the user, nil) is returned. At this point, the user has been validated
// and LoginUser can be called.
//
// If an error occurs, _either_ the record does not exist, or some database issue prevented deletion,
// and in either event, the user cannot be approved. In this case (0, <error>) will be returned
func (ah AShirtAuthBridge) OneTimeVerification(ctx context.Context, userKey string, expirationInMinutes int64) (int64, error) {

	var userID int64
	err := ah.db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Get(&userID, sq.Select("user_id").From("auth_scheme_data").
			Where(sq.Eq{"user_key": userKey}).                                                  // The recovery code exists...
			Where("TIMESTAMPDIFF(minute, created_at, ?) < ?", time.Now(), expirationInMinutes)) // and the record hasn't expired

		tx.Delete(sq.Delete("auth_scheme_data").Where(sq.Eq{"user_key": userKey}))
	})
	if err != nil {
		return 0, backend.WrapError("Unable to validate one-time verification", err)
	}
	return userID, nil
}

// GetDatabase provides raw access to the database. In general, this should not be used by authschemes,
// but is provided in situations where unique-access to the database is required.
func (ah AShirtAuthBridge) GetDatabase() *database.Connection {
	return ah.db
}

// AddScheduledEmail creates a database entry for an outgoing email, for the given email address and
// related user_id
func (ah AShirtAuthBridge) AddScheduledEmail(emailAddress string, userID int64, emailTemplate string) error {
	_, err := ah.db.Insert("email_queue", map[string]interface{}{
		"to_email": emailAddress,
		"user_id":  userID,
		"template": emailTemplate,
	})
	if err != nil {
		return backend.WrapError("Unable to schedule email", backend.DatabaseErr(err))
	}
	return nil
}
