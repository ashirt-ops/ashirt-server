// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package localauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/localauth/constants"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/server/remux"
	"golang.org/x/crypto/bcrypt"
)

// LocalAuthScheme is a small structure capturing the data requirements specific to local authentication
type LocalAuthScheme struct {
	RegistrationEnabled bool
}

// Name returns the name of this authscheme
func (LocalAuthScheme) Name() string {
	return constants.Code
}

// FriendlyName returns "ASHIRT Local Authentication"
func (LocalAuthScheme) FriendlyName() string {
	return constants.FriendlyName
}

// Flags returns auth flags associated with local auth
// in particular, notes if registration is open or closed
func (s LocalAuthScheme) Flags() []string {
	flags := make([]string, 0)

	if s.RegistrationEnabled {
		flags = append(flags, "open-registration")
	}

	return flags
}

func (LocalAuthScheme) Type() string {
	return constants.Code
}

// BindRoutes creates many routes for local database routes:
//
// * POST   ${prefix}/register             Flags that a new user should be created
//
// * POST   ${prefix}/login                Verifies the username/password combo
//
// * POST   ${prefix}/login/resetpassword  Second authentication step for users to reset their password if forced to
//
// * PUT    ${prefix}/password             Allows users to change their password
//
// * PUT    ${prefix}/admin/password       Allows admins to reset a user's password
//
// * POST   ${prefix}/admin/register       Allows admins to create new users on behalf of that user.
//
// * POST   ${prefix}/link                 Adds local auth to a non-local user
//
// * TOTP-Related
//   - POST   ${prefix}/login/totp     Completes login with totp passcode
//   - GET    ${prefix}/totp           Returns boolean true if the user has totp enabled, false otherwise
//   - GET    ${prefix}/totp/generate  Returns a new generated totp secret/uri/qrcode
//   - POST   ${prefix}/totp           Enables totp on a user's account by accepting a secret and verifying
//     a corresponding one time passcode (errors if one already exists)
//   - DELETE ${prefix}/totp           Removes a totp secret from a user's account
//
// In each case above, the actual action is deferred to the bridge connecting this auth scheme to
// the underlying system/database
func (p LocalAuthScheme) BindRoutes(r *mux.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "POST", "/register", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		if !p.RegistrationEnabled {
			return nil, fmt.Errorf("registration is closed to users")
		}

		dr := remux.DissectJSONRequest(r)
		info := RegistrationInfo{
			Username:  dr.FromBody("username").Required().AsString(),
			Email:     dr.FromBody("email").Required().AsString(),
			FirstName: dr.FromBody("firstName").Required().AsString(),
			LastName:  dr.FromBody("lastName").Required().AsString(),
			Password:  dr.FromBody("password").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}

		if err := checkPasswordComplexity(info.Password); err != nil {
			return nil, err
		}
		if err := bridge.ValidateRegistrationInfo(info.Email, info.Username); err != nil {
			return nil, err
		}

		return nil, registerNewUser(r.Context(), bridge, info)
	}))

	remux.Route(r, "POST", "/admin/register", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		if !middleware.IsAdmin(r.Context()) {
			return nil, backend.UnauthorizedWriteErr(fmt.Errorf("Requesting user is not an admin"))
		}

		authKey := make([]byte, 42)
		if _, err := rand.Read(authKey); err != nil {
			return nil, backend.WrapError("Unable to generate random new user password key", err)
		}

		// convert authKey into readable format
		readKey := base64.StdEncoding.EncodeToString(authKey)

		dr := remux.DissectJSONRequest(r)
		info := RegistrationInfo{
			Username:           dr.FromBody("username").Required().AsString(),
			Email:              dr.FromBody("email").Required().AsString(),
			FirstName:          dr.FromBody("firstName").Required().AsString(),
			LastName:           dr.FromBody("lastName").AsString(),
			Password:           string(readKey),
			ForceResetPassword: true,
		}

		if dr.Error != nil {
			return nil, dr.Error
		}

		if err := checkPasswordComplexity(info.Password); err != nil {
			return nil, err
		}

		if err := registerNewUser(r.Context(), bridge, info); err != nil {
			return nil, err
		}

		return dtos.NewUserCreatedByAdmin{
			TemporaryPassword: info.Password,
		}, nil
	}))

	remux.Route(r, "POST", "/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			username := dr.FromBody("username").Required().AsString()
			password := dr.FromBody("password").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}

			authData, findUserErr := bridge.FindUserAuth(username)
			checkPwErr := checkUserPassword(authData, password)
			if firstErr := backend.FirstError(findUserErr, checkPwErr); firstErr != nil {
				return nil, backend.WrapError("Could not validate user", backend.InvalidCredentialsErr(firstErr))
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

			sess := readLocalSession(r, bridge)
			if !sess.SessionValid {
				return nil, backend.HTTPErr(http.StatusUnauthorized,
					"Your account does not require a password reset at this time",
					errors.New("User session is not a local auth needsPasswordResetAuthSession"))
			}

			if err := updateUserPassword(bridge, sess.Username, newPassword); err != nil {
				return nil, backend.WrapError("Unable to reset user password", err)
			}

			authData, err := bridge.FindUserAuth(sess.Username)
			if err != nil {
				return nil, backend.WrapError("Unable to reset user password", err)
			}

			return nil, attemptFinishLogin(w, r, bridge, authData)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/login/totp", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			totpPasscode := dr.FromBody("totpPasscode").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}

			sess := readLocalSession(r, bridge)
			if !sess.SessionValid {
				return nil, backend.HTTPErr(http.StatusUnauthorized,
					"Could not validate passcode",
					errors.New("User session does not require needsTotpAuthSession"))
			}

			authData, err := bridge.FindUserAuth(sess.Username)

			if authData.TOTPSecret == nil {
				return nil, backend.HTTPErr(http.StatusUnauthorized,
					"Could not validate passcode",
					errors.New("User trying to authenticate with TOTP when TOTP is not enabled"))
			}

			if err = validateTOTP(totpPasscode, *authData.TOTPSecret); err != nil {
				return nil, backend.WrapError("Could not validate passcode", err)
			}
			sess.TOTPValidated = true
			if err = sess.writeLocalSession(w, r, bridge); err != nil {
				return nil, backend.WrapError("Could not validate passcode", backend.WrapError("Unable to set auth scheme in session", err))
			}

			return nil, attemptFinishLogin(w, r, bridge, authData)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "PUT", "/password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			username := dr.FromBody("username").Required().AsString()
			oldPassword := dr.FromBody("oldPassword").Required().AsString()
			newPassword := dr.FromBody("newPassword").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}

			authData, findUserErr := bridge.FindUserAuth(username)
			checkPwErr := checkUserPassword(authData, oldPassword)
			if firstErr := backend.FirstError(findUserErr, checkPwErr); firstErr != nil {
				return nil, backend.WrapError("Unable to set new password", backend.InvalidPasswordErr(firstErr))
			}
			if authData.UserID != middleware.UserID(r.Context()) {
				return nil, backend.InvalidPasswordErr(errors.New("Cannot reset password for a different user than is currently logged in"))
			}

			return nil, updateUserPassword(bridge, username, newPassword)
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

		// TODO admin reset should be providing username instead of userSlug and this method should be deleted from auth bridge:
		userAuths, err := bridge.FindUserAuthsByUserSlug(userSlug)
		if err != nil {
			return nil, err
		}
		if len(userAuths) != 1 {
			return nil, fmt.Errorf("More than one local auth row exists for user %s", userSlug)
		}
		userAuth := userAuths[0]

		_, err = bridge.FindUserAuth(userAuth.Username)

		if err != nil {
			return nil, backend.NotFoundErr(fmt.Errorf("User %v does not have %v authentication", userSlug, p.Name()))
		}

		// Skipping password requirement check here -- Admins should have free reign
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, backend.WrapError("Unable to encrypt new password", err)
		}

		userAuth.EncryptedPassword = encryptedPassword
		userAuth.NeedsPasswordReset = true

		return nil, bridge.UpdateAuthForUser(userAuth)
	}))

	remux.Route(r, "POST", "/link", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		username := dr.FromBody("username").Required().AsString()
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

		callingUserID := middleware.UserID(r.Context())
		if err := bridge.ValidateLinkingInfo(username, callingUserID); err != nil {
			return nil, err
		}

		err = bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:            callingUserID,
			Username:          username,
			EncryptedPassword: encryptedPassword,
		})

		return nil, err
	}))

	remux.Route(r, "GET", "/totp", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		userSlug := dr.FromBody("userSlug").AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return readUserTotpStatus(r.Context(), bridge, userSlug)
	}))

	remux.Route(r, "GET", "/totp/generate", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		userAuth, err := bridge.FindUserAuthByContext(r.Context())
		if err != nil {
			return nil, err
		}

		return generateTOTP(userAuth.Username)
	}))

	remux.Route(r, "POST", "/totp", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		secret := dr.FromBody("secret").Required().AsString()
		passcode := dr.FromBody("passcode").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}

		userAuth, err := bridge.FindUserAuthByContext(r.Context())
		if err != nil {
			return nil, err
		}
		if userAuth.TOTPSecret != nil {
			return nil, backend.BadInputErr(
				errors.New("User already has a TOTP key associated"),
				"Your account already has a TOTP key",
			)
		}

		err = validateTOTP(passcode, secret)
		if err != nil {
			return nil, err
		}

		userAuth.TOTPSecret = &secret
		err = bridge.UpdateAuthForUser(userAuth)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}))

	remux.Route(r, "DELETE", "/totp", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		userSlug := dr.FromBody("userSlug").AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}

		return nil, deleteUserTotp(r.Context(), bridge, userSlug)
	}))
}

func attemptFinishLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, authData authschemes.UserAuthData) error {
	sess := readLocalSession(r, bridge)
	sess.Username = authData.Username

	if authData.TOTPSecret != nil {
		if !sess.SessionValid || !sess.TOTPValidated {
			sess.TOTPValidated = false
			if err := sess.writeLocalSession(w, r, bridge); err != nil {
				return backend.WrapError("Unable to set auth scheme in session", err)
			}
			return backend.UserRequiresAdditionalAuthenticationErr("TOTP_REQUIRED")
		}
	}

	if authData.NeedsPasswordReset {
		if err := sess.writeLocalSession(w, r, bridge); err != nil {
			return backend.WrapError("Unable to set auth scheme in session", err)
		}
		return backend.UserRequiresAdditionalAuthenticationErr("PASSWORD_RESET_REQUIRED")
	}

	if err := bridge.LoginUser(w, r, authData.UserID, nil); err != nil {
		return backend.WrapError("Attempt to finish login failed", err)
	}

	return nil
}

func updateUserPassword(bridge authschemes.AShirtAuthBridge, username string, newPassword string) error {
	authData, err := bridge.FindUserAuth(username)
	if err != nil {
		return backend.WrapError("Unable to update password", err)
	}

	if err = checkPasswordComplexity(newPassword); err != nil {
		return backend.WrapError("Unable to update password", err)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return backend.WrapError("Unable to encrypte new password", err)
	}

	authData.EncryptedPassword = encryptedPassword
	authData.NeedsPasswordReset = false

	return bridge.UpdateAuthForUser(authData)
}

func checkUserPassword(authData authschemes.UserAuthData, password string) error {
	return bcrypt.CompareHashAndPassword(authData.EncryptedPassword, []byte(password))
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
