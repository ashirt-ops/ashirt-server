// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package localauth

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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
// * PUT    ${prefix}/password/admin       Allows admins to reset a user's password
//
// * POST   ${prefix}/link                 Adds local auth to a non-local user
//
// * TOTP-Related
//   * GET    ${prefix}/totp           Returns boolean true if the user has totp enabled, false otherwise
//   * GET    ${prefix}/totp/generate  Returns a new generated totp secret/uri/qrcode
//   * POST   ${prefix}/totp           Enables totp on a user's account by accepting a secret and verifying
//     a corresponding one time passcode (errors if one already exists)
//   * DELETE ${prefix}/totp           Removes a totp secret from a user's account
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

			authData, findUserErr := bridge.FindUserAuth(userKey)
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

			sess, ok := bridge.ReadAuthSchemeSession(r).(*needsPasswordResetAuthSession)
			if !ok {
				return nil, backend.HTTPErr(http.StatusUnauthorized,
					"Your account does not require a password reset at this time",
					errors.New("User session is not a local auth needsPasswordResetAuthSession"))
			}

			if err := updateUserPassword(bridge, sess.UserKey, newPassword); err != nil {
				return nil, backend.WrapError("Unable to reset user password", err)
			}

			authData, err := bridge.FindUserAuth(sess.UserKey)
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

			sess, ok := bridge.ReadAuthSchemeSession(r).(*needsPasswordResetAuthSession)
			if !ok {
				return nil, backend.HTTPErr(http.StatusUnauthorized,
					"Could not validate passcode",
					errors.New("User session does not require needsTotpAuthSession"))
			}

			authData, err := bridge.FindUserAuth(sess.UserKey)

			if authData.TOTPSecret == nil {
				return nil, backend.HTTPErr(http.StatusUnauthorized,
					"Could not validate passcode",
					errors.New("User trying to authenticate with TOTP when TOTP is not enabled"))
			}

			if err = validateTOTP(totpPasscode, *authData.TOTPSecret); err != nil {
				return nil, backend.WrapError("Could not validate passcode", err)
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

			authData, findUserErr := bridge.FindUserAuth(userKey)
			checkPwErr := checkUserPassword(authData, oldPassword)
			if firstErr := backend.FirstError(findUserErr, checkPwErr); firstErr != nil {
				return nil, backend.WrapError("Unable to set new password", backend.InvalidPasswordErr(firstErr))
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
		userAuths, err := bridge.FindUserAuthsByUserSlug(userSlug)
		if err != nil {
			return nil, err
		}
		if len(userAuths) != 1 {
			return nil, fmt.Errorf("More than one local auth row exists for user %s", userSlug)
		}
		userAuth := userAuths[0]

		_, err = bridge.FindUserAuth(userAuth.UserKey)

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

	remux.Route(r, "GET", "/totp", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		userAuth, err := bridge.FindUserAuthByContext(r.Context())
		if err != nil {
			return nil, err
		}
		hasTOTP := userAuth.TOTPSecret != nil
		return hasTOTP, nil
	}))

	remux.Route(r, "GET", "/totp/generate", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		userAuth, err := bridge.FindUserAuthByContext(r.Context())
		if err != nil {
			return nil, err
		}

		return generateTOTP(userAuth.UserKey)
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

		fmt.Println(secret, passcode) // TODO: remove
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
		userAuth, err := bridge.FindUserAuthByContext(r.Context())
		if err != nil {
			return nil, err
		}
		if userAuth.TOTPSecret == nil {
			return nil, backend.BadInputErr(
				errors.New("User does not have a TOTP key associated"),
				"Your account does not have a TOTP key",
			)
		}

		userAuth.TOTPSecret = nil
		return nil, bridge.UpdateAuthForUser(userAuth)
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
	if authData.TOTPSecret != nil {
		err := bridge.SetAuthSchemeSession(w, r, &needsTotpAuthSession{UserKey: authData.UserKey})
		if err != nil {
			return backend.WrapError("Unable to set auth scheme in session", err)
		}
		return backend.UserRequiresAdditionalAuthenticationErr("TOTP_REQUIRED")
	}

	err := bridge.LoginUser(w, r, authData.UserID, nil)
	if err != nil {
		return backend.WrapError("Attempt to finish login failed", err)
	}

	return nil
}

func updateUserPassword(bridge authschemes.AShirtAuthBridge, userKey string, newPassword string) error {
	authData, err := bridge.FindUserAuth(userKey)
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

// TOTPKey represents the secret and QR code for a given URL to authenticate with
type TOTPKey struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
	QRCode string `json:"qr"`
}

const totpDigits = otp.DigitsSix
const totpAlgorithm = otp.AlgorithmSHA1

func generateTOTP(accountName string) (*TOTPKey, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ashirt",
		AccountName: accountName,
		SecretSize:  64,
		Digits:      totpDigits,
		Algorithm:   totpAlgorithm,
	})
	if err != nil {
		return nil, err
	}
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)

	return &TOTPKey{
		URL:    key.URL(),
		Secret: key.Secret(),
		QRCode: "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

func validateTOTP(passcode string, totpSecret string) error {
	if passcode == "" {
		return backend.MissingValueErr("TOTP Passcode")
	}
	if totpSecret == "" {
		return backend.MissingValueErr("TOTP Secret")
	}

	isValid, err := totp.ValidateCustom(passcode, totpSecret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    totpDigits,
		Algorithm: totpAlgorithm,
	})
	if err != nil {
		return backend.InvalidTOTPErr(err)
	}
	if !isValid {
		return backend.InvalidTOTPErr(errors.New("totp.Validate failure"))
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
