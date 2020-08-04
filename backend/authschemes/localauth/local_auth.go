// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package localauth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/localauth/constants"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/server/remux"
	"golang.org/x/crypto/bcrypt"
)

type LocalAuthScheme struct{}

// Name returns the name of this authscheme
func (LocalAuthScheme) Name() string {
	return constants.Code
}

// FriendlyName returns "ASHIRT Local Authentication"
func (LocalAuthScheme) FriendlyName() string {
	return constants.FriendlyName
}

// BindRoutes creates four routes for local database routes:
//
// 1. POST ${prefix}/register (flags that a new user should be created)
//
// 2. POST ${prefix}/login (verifies the username/password combo)
//
// 3. POST ${prefix}/login/resetpassword (second authentication step for users to reset their password if forced to)
//
// 4. PUT ${prefix}/password (allows users to change their password)
//
// 5. PUT ${prefix}/password/admin (allows admins to reset a user's password)
//
// 6. POST ${prefix}/link (adds local auth to a non-local user)
//
// In each case above, the actual action is deferred to the bridge connecting this auth scheme to
// the underlying system/database
func (p LocalAuthScheme) BindRoutes(r *mux.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "POST", "/register", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		firstName := dr.FromBody("firstName").Required().AsString()
		lastName := dr.FromBody("lastName").Required().AsString()
		email := dr.FromBody("email").Required().AsString()
		password := dr.FromBody("password").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}

		if err := checkPasswordComplexity(password); err != nil {
			return nil, err
		}

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, backend.WrapError("Unable to generate encrypted password", err)
		}

		userResult, err := bridge.CreateNewUser(authschemes.UserProfile{
			FirstName: firstName,
			LastName:  lastName,
			Slug:      strings.ToLower(firstName + "." + lastName),
			Email:     email,
		})
		if err != nil {
			return nil, err
		}

		return nil, bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:            userResult.UserID,
			UserKey:           email,
			EncryptedPassword: encryptedPassword,
		})
	}))

	remux.Route(r, "POST", "/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			userKey := dr.FromBody("email").Required().AsString()
			password := dr.FromBody("password").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}

			authData, authDataErr := bridge.FindUserAuth(userKey)
			checkPwErr := checkUserCredentials(authData, password)
			if authDataErr != nil || checkPwErr != nil {
				return nil, backend.WrapError("Could not validate user", backend.InvalidCredentialsErr(coalesceError(authDataErr, checkPwErr)))
			}

			return nil, attemptFinishLogin(w, r, bridge, authData)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/login/resetpassword", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			newPassword := dr.FromBody("newPassword").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}

			sess, ok := bridge.ReadAuthSchemeSession(r).(*needsPasswordResetAuthSession)
			if !ok {
				return nil, backend.HTTPErr(http.StatusUnauthorized, "Your account does not require a password reset at this time", errors.New("User session is not a local auth needsPasswordResetAuthSession"))
			}

			err := updateUserPassword(bridge, sess.UserKey, newPassword)
			if err != nil {
				return nil, backend.WrapError("Unable to reset user password", err)
			}

			authData, err := bridge.FindUserAuth(sess.UserKey)
			if err != nil {
				return nil, backend.WrapError("Unable to reset user password", err)
			}

			return nil, attemptFinishLogin(w, r, bridge, authData)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "PUT", "/password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			userKey := dr.FromBody("userKey").Required().AsString()
			oldPassword := dr.FromBody("oldPassword").Required().AsString()
			newPassword := dr.FromBody("newPassword").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}

			authData, authDataErr := bridge.FindUserAuth(userKey)
			checkPwErr := checkUserCredentials(authData, oldPassword)
			if authDataErr != nil || checkPwErr != nil {
				return nil, backend.WrapError("Unable to set new password", backend.InvalidPasswordErr(coalesceError(authDataErr, checkPwErr)))
			}
			if authData.UserID != middleware.UserID(r.Context()) {
				return nil, backend.InvalidPasswordErr(errors.New("Cannot reset password for a different user than is currently logged in"))
			}

			return nil, updateUserPassword(bridge, userKey, newPassword)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "PUT", "/admin/password", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		if !middleware.IsAdmin(r.Context()) {
			return nil, backend.UnauthorizedWriteErr(fmt.Errorf("Requesting user is not an admin"))
		}

		dr := remux.DissectJSONRequest(r)
		userSlug := dr.FromBody("userSlug").Required().AsString()
		newPassword := dr.FromBody("newPassword").Required().AsString()

		if dr.Error != nil {
			return nil, dr.Error
		}

		// TODO admin reset should be providing userKey instead of userSlug and this method should be deleted from auth bridge:
		profiles, err := bridge.FindUserAuthsByUserSlug(userSlug)
		if err != nil {
			return nil, err
		}
		if len(profiles) != 1 {
			return nil, fmt.Errorf("More than one local auth row exists for user %s", userSlug)
		}
		profile := profiles[0]

		_, err = bridge.FindUserAuth(profile.UserKey)

		if err != nil {
			return nil, backend.NotFoundErr(fmt.Errorf("User %v does not have %v authentication", userSlug, p.Name()))
		}

		// Skipping password requirement check here -- Admins should have free reign
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, backend.WrapError("Unable to encrypt new password", err)
		}

		return nil, bridge.UpdateAuthForUser(profile.UserKey, encryptedPassword, true)
	}))

	remux.Route(r, "POST", "/link", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		email := dr.FromBody("email").Required().AsString()
		password := dr.FromBody("password").Required().AsString()

		if dr.Error != nil {
			return nil, dr.Error
		}

		if err := checkPasswordComplexity(password); err != nil {
			return nil, err
		}

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, backend.WrapError("Unable to encrypt new password", err)
		}

		err = bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:            middleware.UserID(r.Context()),
			UserKey:           email,
			EncryptedPassword: encryptedPassword,
		})

		return nil, err
	}))
}

func attemptFinishLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, authData authschemes.UserAuthData) error {
	if authData.NeedsPasswordReset {
		err := bridge.SetAuthSchemeSession(w, r, &needsPasswordResetAuthSession{UserKey: authData.UserKey})
		if err != nil {
			return backend.WrapError("Unable to set auth scheme in session", err)
		}
		return backend.UserRequiresAdditionalAuthenticationErr("PASSWORD_RESET_REQUIRED")
	}

	err := bridge.LoginUser(w, r, authData.UserID, nil)
	if err != nil {
		return backend.WrapError("Attempt to finish login failed", err)
	}

	return nil
}

func updateUserPassword(bridge authschemes.AShirtAuthBridge, userKey string, newPassword string) error {
	if err := checkPasswordComplexity(newPassword); err != nil {
		return backend.WrapError("Unable to update password", err)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return backend.WrapError("Unable to encrypte new password", err)
	}

	return bridge.UpdateAuthForUser(userKey, encryptedPassword, false)
}

func checkUserCredentials(authData authschemes.UserAuthData, password string) error {
	return bcrypt.CompareHashAndPassword(authData.EncryptedPassword, []byte(password))
}

// coalesceError returns the first non-nil error. Will return the following: e1, then e2, then nil
func coalesceError(e1 error, e2 error) error {
	if e1 != nil {
		return e1
	}
	if e2 != nil {
		return e2
	}
	return nil
}

func checkPasswordComplexity(suggestedPassword string) error {
	err := errors.New("Password did not meet requirements")
	if len(suggestedPassword) < 5 {
		return backend.BadInputErr(err, "Password must be at least 5 characters long")
	}
	// TODO: Fill in with password complexity requirements/tests
	// if strings.Contains(suggestedPassword, "a") {
	//   return backend.BadInputErr(err, "Password must not use an `a` character")
	// }
	// if !strings.Contains(suggestedPassword, "0") {
	//   return backend.BadInputErr(err, "Password must contain a 0 character")
	// }

	return nil
}
