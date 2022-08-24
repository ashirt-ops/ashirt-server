// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package webauthn

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/webauthn/constants"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/server/remux"

	"github.com/duo-labs/webauthn/protocol"
	auth "github.com/duo-labs/webauthn/webauthn"
)

type WebAuthn struct {
	RegistrationEnabled bool
	Web                 *auth.WebAuthn
}

func New(cfg config.AuthInstanceConfig, webConfig *config.WebConfig) (WebAuthn, error) {
	parsedUrl, err := url.Parse(webConfig.FrontendIndexURL)

	var host string
	var port string
	if err != nil {
		return WebAuthn{}, err
	}

	if host, port, err = net.SplitHostPort(parsedUrl.Host); err != nil {
		host = parsedUrl.Host
	}

	// TODO: I don't understand how to correctly set the RPOrigin. the code works *specifically* for
	// localhost, but may fail for proper deployments. We might need to make this an env var.
	config := auth.Config{
		RPDisplayName: cfg.DisplayName,
		RPID:          host,
	}

	if host == "localhost" {
		config.RPOrigin = "http://" + host + ":" + port
	}

	web, err := auth.New(&config)
	if err != nil {
		return WebAuthn{}, err
	}

	return WebAuthn{
		RegistrationEnabled: cfg.RegistrationEnabled,
		Web:                 web,
	}, nil
}

func (a WebAuthn) Name() string {
	return constants.Name
}

func (a WebAuthn) FriendlyName() string {
	return constants.FriendlyName
}

func (a WebAuthn) Flags() []string {
	return []string{}
}

func (a WebAuthn) Type() string {
	return constants.Name
}

func (a WebAuthn) BindRoutes(r *mux.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "POST", "/register/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			// validate basic registration data
			if !a.RegistrationEnabled {
				return nil, errors.New("registration is closed to users")
			}

			dr := remux.DissectJSONRequest(r)
			info := WebAuthnRegistrationInfo{
				Email:     dr.FromBody("email").Required().AsString(),
				FirstName: dr.FromBody("firstName").Required().AsString(),
				LastName:  dr.FromBody("lastName").Required().AsString(),
				KeyName:   dr.FromBody("keyName").Required().AsString(),
			}
			if dr.Error != nil {
				return nil, dr.Error
			}

			if taken, err := bridge.CheckIfUserEmailTaken(info.Email, -1, true); err != nil {
				return nil, backend.BadAuthErr(errors.New("Unable to review user"))
			} else if taken {
				return nil, backend.BadAuthErr(errors.New("User has already been registered. If you are this user, please link your account instead"))
			}

			return a.beginRegistration(w, r, bridge, info)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/register/finish", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		data, encodedCreds, err := a.validateRegistrationComplete(r, bridge)
		if err != nil {
			return nil, backend.WrapError("Unable to validate registration data", err)
		}

		userProfile := authschemes.UserProfile{
			FirstName: data.UserData.FirstName,
			LastName:  data.UserData.LastName,
			Slug:      data.UserData.Email,
			Email:     data.UserData.Email,
		}
		userResult, err := bridge.CreateNewUser(userProfile)
		if err != nil {
			return nil, backend.WrapError("Unable to create user", err)
		}

		return nil, bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:   userResult.UserID,
			UserKey:  userProfile.Email,
			JSONData: helpers.Ptr(string(encodedCreds)),
		})
	}))

	remux.Route(r, "POST", "/login/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			email := dr.FromBody("email").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}
			return a.beginLogin(w, r, bridge, email)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/login/finish", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			rawData := bridge.ReadAuthSchemeSession(r)
			data, ok := rawData.(*webAuthNSessionData)
			if !ok {
				return nil, errors.New("Unable to complete login -- session not found or corrupt")
			}

			user := &data.UserData
			cred, err := a.Web.FinishLogin(user, *data.WebAuthNSessionData, r)
			if err != nil {
				return nil, backend.BadAuthErr(err)
			} else if cred.Authenticator.CloneWarning {
				return nil, backend.WrapError("credential appears to be cloned", backend.BadAuthErr(err))
			}

			updateSignCount(data, cred, bridge)

			if err := bridge.LoginUser(w, r, data.UserData.UserIDAsI64(), nil); err != nil {
				return nil, backend.WrapError("Attempt to finish login failed", err)
			}

			return nil, nil
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/link/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			callingUserId := middleware.UserID(r.Context())

			dr := remux.DissectJSONRequest(r)
			info := WebAuthnRegistrationInfo{
				Email:   dr.FromBody("email").Required().AsString(),
				KeyName: dr.FromBody("keyName").Required().AsString(),
				UserID:  callingUserId,
			}
			if dr.Error != nil {
				return nil, dr.Error
			}

			if emailTaken, err := bridge.CheckIfUserEmailTaken(info.Email, callingUserId, true); err != nil {
				return nil, err
			} else if emailTaken {
				return nil, backend.BadInputErr(
					errors.New("error linking account: email taken"),
					"An account for this user already exists",
				)
			}

			return a.beginRegistration(w, r, bridge, info)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/link/finish", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		data, encodedCreds, err := a.validateRegistrationComplete(r, bridge)
		if err != nil {
			return nil, backend.WrapError("Unable to validate registration data", err)
		}
		return nil, bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:   byteSliceToI64(data.UserData.UserID),
			UserKey:  data.UserData.Email,
			JSONData: helpers.Ptr(string(encodedCreds)),
		})
	}))

	remux.Route(r, "GET", "/keys", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		callingUserID := middleware.UserID(r.Context())
		return a.getKeys(callingUserID, bridge)
	}))

	remux.Route(r, "DELETE", "/key/{keyName}", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		callingUserID := middleware.UserID(r.Context())
		dr := remux.DissectJSONRequest(r)
		keyName := dr.FromURL("keyName").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, a.deleteKey(callingUserID, keyName, bridge)
	}))

	remux.Route(r, "POST", "/key/add/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			callingUserId := middleware.UserID(r.Context())

			auth, err := bridge.FindUserAuthByUserID(callingUserId)
			if err != nil {
				return nil, backend.DatabaseErr(err)
			}

			dr := remux.DissectJSONRequest(r)
			info := WebAuthnRegistrationInfo{
				Email:   auth.UserKey,
				KeyName: dr.FromBody("keyName").Required().AsString(),
				UserID:  callingUserId,
			}
			if dr.Error != nil {
				return nil, dr.Error
			}

			if emailTaken, err := bridge.CheckIfUserEmailTaken(info.Email, callingUserId, true); err != nil {
				return nil, err
			} else if emailTaken {
				return nil, backend.BadInputErr(
					errors.New("error linking account: email taken"),
					"An account for this user already exists",
				)
			}

			return a.beginRegistration(w, r, bridge, info)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/key/add/finish", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		_, encodedCreds, err := a.validateRegistrationComplete(r, bridge)
		if err != nil {
			return nil, backend.WrapError("Unable to validate registration data", err)
		}

		callingUserID := middleware.UserID(r.Context())
		userAuth, err := bridge.FindUserAuthByUserID(callingUserID)
		if err != nil {
			return nil, backend.WrapError("Unable to find user", err)
		}
		userAuth.JSONData = helpers.Ptr(string(encodedCreds))
		err = bridge.UpdateAuthForUser(userAuth)
		if err != nil {
			return nil, backend.WrapError("Unable to update keys", err)
		}

		// We might want to return a full list of keys. TODO: check if we want that
		return nil, nil
	}))
}

func (a WebAuthn) getKeys(userID int64, bridge authschemes.AShirtAuthBridge) (*ListKeysOutput, error) {
	auth, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return nil, backend.WrapError("Unable to get keys", err)
	}

	webauthRawCreds := []byte(*auth.JSONData)
	var creds []AShirtWebauthnCredential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return nil, backend.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	results := helpers.Map(creds, func(cred AShirtWebauthnCredential) string {
		return cred.KeyName
	})
	output := ListKeysOutput{results}
	return &output, nil
}

func (a WebAuthn) deleteKey(userID int64, keyName string, bridge authschemes.AShirtAuthBridge) error {
	auth, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return backend.WrapError("Unable to find user", err)
	}

	webauthRawCreds := []byte(*auth.JSONData)
	var creds []AShirtWebauthnCredential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return backend.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	results := helpers.Filter(creds, func(cred AShirtWebauthnCredential) bool {
		return cred.KeyName != keyName
	})
	encodedCreds, err := json.Marshal(results)
	if err != nil {
		return backend.WrapError("Unable to delete key", err)
	}
	auth.JSONData = helpers.Ptr(string(encodedCreds))

	bridge.UpdateAuthForUser(auth)

	return nil
}

func (a WebAuthn) beginRegistration(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, info WebAuthnRegistrationInfo) (*protocol.CredentialCreation, error) {
	var user webauthnUser
	if info.UserID == 0 {
		user = makeNewWebAuthnUser(info.FirstName, info.LastName, info.Email, info.KeyName)
	} else {
		user = makeLinkingWebAuthnUser(info.UserID, info.Email, info.KeyName)
		authData, err := bridge.FindUserAuthByUserID(info.UserID)
		if err != nil {
			return nil, backend.WebauthnLoginError(err, "Unable to find existing user auth")
		}
		creds, err := a.getExistingCredentials(authData)
		if err != nil {
			return nil, backend.WebauthnLoginError(err, "Unable to parse webauthn credentials")
		}
		user.Credentials = append(user.Credentials, creds...)
	}

	credExcludeList := make([]protocol.CredentialDescriptor, len(user.Credentials))
	for i, cred := range user.Credentials {
		credExcludeList[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
	}
	registrationOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = credExcludeList
	}

	credOptions, sessionData, err := a.Web.BeginRegistration(&user, registrationOptions)
	if err != nil {
		return nil, err
	}

	err = bridge.SetAuthSchemeSession(w, r, makeWebauthNSessionData(user, sessionData))

	return credOptions, err
}

func (a WebAuthn) beginLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, email string) (interface{}, error) {
	authData, err := bridge.FindUserAuth(email)
	if err != nil {
		return nil, backend.WebauthnLoginError(err, "Could not validate user", "No such auth")
	}
	if authData.JSONData == nil {
		return nil, backend.WebauthnLoginError(err, "User lacks webauthn credentials")
	}

	user, err := bridge.GetUserFromID(authData.UserID)
	if err != nil {
		return nil, backend.WebauthnLoginError(err, "Could not validate user", "No such user")
	}

	creds, err := a.getExistingCredentials(authData)
	if err != nil {
		return nil, backend.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	webauthnUser := makeWebAuthnUser(user.FirstName, user.LastName, user.Slug, user.Email, user.ID, creds)
	options, sessionData, err := a.Web.BeginLogin(&webauthnUser)
	if err != nil {
		return nil, backend.WebauthnLoginError(err, "Unable to begin login process")
	}

	if err = bridge.SetAuthSchemeSession(w, r, makeWebauthNSessionData(webauthnUser, sessionData)); err != nil {
		return nil, backend.WebauthnLoginError(err, "Unable to begin login process", "Unable to set session")
	}

	return options, nil
}

func (a WebAuthn) getExistingCredentials(authData authschemes.UserAuthData) ([]AShirtWebauthnCredential, error) {
	webauthRawCreds := []byte(*authData.JSONData)
	var creds []AShirtWebauthnCredential
	if err := json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return nil, backend.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	return creds, nil
}

func (a WebAuthn) validateRegistrationComplete(r *http.Request, bridge authschemes.AShirtAuthBridge) (*webAuthNSessionData, []byte, error) {
	rawData := bridge.ReadAuthSchemeSession(r)
	data, ok := rawData.(*webAuthNSessionData)
	if !ok {
		return nil, nil, errors.New("Unable to complete registration -- session not found or corrupt")
	}

	cred, err := a.Web.FinishRegistration(&data.UserData, *data.WebAuthNSessionData, r)
	if err != nil {
		return nil, nil, backend.WrapError("Unable to complete registration", err)
	}

	data.UserData.Credentials = append(data.UserData.Credentials, wrapCredential(*cred, data.UserData.KeyName))

	encodedCreds, err := json.Marshal(data.UserData.Credentials)
	if err != nil {
		return nil, nil, backend.WrapError("Unable to create registration", err)
	}
	return data, encodedCreds, nil
}

func updateSignCount(data *webAuthNSessionData, loginCred *auth.Credential, bridge authschemes.AShirtAuthBridge) error {
	userID := data.UserData.UserIDAsI64()

	userAuth, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return backend.WrapError("Unable to find user", err)
	}
	matchingIndex, _ := helpers.Find(data.UserData.Credentials, func(item AShirtWebauthnCredential) bool {
		return string(item.ID) == string(loginCred.ID)
	})
	if matchingIndex == -1 {
		return backend.WrapError("Could not find matching credential", err)
	}
	data.UserData.Credentials[matchingIndex].Authenticator.SignCount = loginCred.Authenticator.SignCount
	encodedCreds, err := json.Marshal(data.UserData.Credentials)
	if err != nil {
		return backend.WrapError("Unable to encode credentials", err)
	}
	userAuth.JSONData = helpers.Ptr(string(encodedCreds))
	if err = bridge.UpdateAuthForUser(userAuth); err != nil {
		return backend.WrapError("Unable to update credential", err)
	}
	return nil
}
