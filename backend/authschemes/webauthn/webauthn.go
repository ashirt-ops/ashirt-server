// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package webauthn

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/webauthn/constants"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/server/remux"

	"github.com/go-webauthn/webauthn/protocol"
	auth "github.com/go-webauthn/webauthn/webauthn"
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

	config := auth.Config{
		RPDisplayName: cfg.WebauthnConfig.DisplayName,
		RPID:          host,
		// the below are all optional
		Debug:                  cfg.WebauthnConfig.Debug,
		Timeout:                cfg.WebauthnConfig.Timeout,
		AttestationPreference:  cfg.WebauthnConfig.Conveyance(),
		AuthenticatorSelection: cfg.BuildAuthenticatorSelection(),
	}

	// TODO: I don't understand how to correctly set the RPOrigin. the code works *specifically* for
	// localhost, but may fail for proper deployments. We might need to make this an env var.
	if cfg.WebauthnConfig.RPOrigin != "" {
		config.RPOrigin = cfg.WebauthnConfig.RPOrigin
	} else if host == "localhost" {
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
	flags := make([]string, 0)

	if a.RegistrationEnabled {
		flags = append(flags, "open-registration")
	}

	return flags
}

func (a WebAuthn) Type() string {
	return constants.Name
}

func (a WebAuthn) BindRoutes(r chi.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "POST", "/register/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			// validate basic registration data
			if !a.RegistrationEnabled {
				return nil, errors.New("registration is closed to users")
			}

			dr := remux.DissectJSONRequest(r)
			info := WebAuthnRegistrationInfo{
				Email:            dr.FromBody("email").Required().AsString(),
				Username:         dr.FromBody("username").Required().AsString(),
				FirstName:        dr.FromBody("firstName").Required().AsString(),
				LastName:         dr.FromBody("lastName").Required().AsString(),
				CredentialName:   dr.FromBody("credentialName").Required().AsString(),
				RegistrationType: CreateCredential,
			}
			if dr.Error != nil {
				return nil, dr.Error
			}

			if err := bridge.ValidateRegistrationInfo(info.Email, info.Username); err != nil {
				return nil, err
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
			Slug:      strings.ToLower(data.UserData.FirstName + "." + data.UserData.LastName),
			Email:     data.UserData.Email,
		}
		userResult, err := bridge.CreateNewUser(userProfile)
		if err != nil {
			return nil, backend.WrapError("Unable to create user", err)
		}

		rawSessionData := bridge.ReadAuthSchemeSession(r)
		sessionData, _ := rawSessionData.(*webAuthNSessionData)

		return nil, bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:   userResult.UserID,
			AuthnID:  sessionData.UserData.AuthnID,
			Username: data.UserData.UserName,
			JSONData: helpers.Ptr(string(encodedCreds)),
		})
	}))

	remux.Route(r, "POST", "/login/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			dr := remux.DissectJSONRequest(r)
			username := dr.FromBody("username").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}
			return a.beginLogin(w, r, bridge, username)
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
				Username:         dr.FromBody("username").Required().AsString(),
				CredentialName:   dr.FromBody("credentialName").Required().AsString(),
				UserID:           callingUserId,
				RegistrationType: LinkCredential,
			}
			if dr.Error != nil {
				return nil, dr.Error
			}

			if err := bridge.ValidateLinkingInfo(info.Username, callingUserId); err != nil {
				return nil, err
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
			Username: data.UserData.UserName,
			JSONData: helpers.Ptr(string(encodedCreds)),
		})
	}))

	remux.Route(r, "GET", "/credentials", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		callingUserID := middleware.UserID(r.Context())
		return a.getCredentials(callingUserID, bridge)
	}))

	remux.Route(r, "DELETE", "/credential/{credentialName}", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		callingUserID := middleware.UserID(r.Context())
		dr := remux.DissectJSONRequest(r)
		credentialName := dr.FromURL("credentialName").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, a.deleteCredential(callingUserID, credentialName, bridge)
	}))

	remux.Route(r, "PUT", "/credential", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		callingUserID := middleware.UserID(r.Context())
		dr := remux.DissectJSONRequest(r)
		if dr.Error != nil {
			return nil, dr.Error
		}
		info := WebAuthnUpdateCredentialInfo{
			NewCredentialName: dr.FromBody("newCredentialName").Required().AsString(),
			CredentialName:    dr.FromBody("credentialName").Required().AsString(),
			UserID:            callingUserID,
		}
		return nil, a.updateCredentialName(info, bridge)
	}))

	remux.Route(r, "POST", "/credential/add/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			auth, err := bridge.FindUserAuthByContext(r.Context())
			if err != nil {
				return nil, backend.DatabaseErr(err)
			}

			dr := remux.DissectJSONRequest(r)
			credentialName := dr.FromBody("credentialName").Required().AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}

			info := WebAuthnRegistrationInfo{
				Username:         auth.Username,
				CredentialName:   credentialName,
				UserID:           auth.UserID,
				RegistrationType: AddCredential,
			}

			creds, err := a.getExistingCredentials(auth)
			if err != nil {
				return nil, err
			}
			info.ExistingCredentials = creds

			return a.beginRegistration(w, r, bridge, info)
		}).ServeHTTP(w, r)
	}))

	remux.Route(r, "POST", "/credential/add/finish", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		_, encodedCreds, err := a.validateRegistrationComplete(r, bridge)
		if err != nil {
			return nil, backend.WrapError("Unable to validate registration data", err)
		}

		userAuth, err := bridge.FindUserAuthByContext(r.Context())
		if err != nil {
			return nil, backend.WrapError("Unable to find user", err)
		}
		userAuth.JSONData = helpers.Ptr(string(encodedCreds))
		err = bridge.UpdateAuthForUser(userAuth)
		if err != nil {
			return nil, backend.WrapError("Unable to update credentials", err)
		}

		// We might want to return a full list of credentials. TODO: check if we want that
		return nil, nil
	}))
}

func (a WebAuthn) getCredentials(userID int64, bridge authschemes.AShirtAuthBridge) (*ListCredentialsOutput, error) {
	auth, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return nil, backend.WrapError("Unable to get credentials", err)
	}

	webauthRawCreds := []byte(*auth.JSONData)
	var creds []AShirtWebauthnCredential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return nil, backend.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	results := helpers.Map(creds, func(cred AShirtWebauthnCredential) CredentialEntry {
		return CredentialEntry{
			CredentialName: cred.CredentialName,
			DateCreated:    cred.CredentialCreatedDate,
		}
	})
	output := ListCredentialsOutput{results}
	return &output, nil
}

func (a WebAuthn) deleteCredential(userID int64, credentialName string, bridge authschemes.AShirtAuthBridge) error {
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
		return cred.CredentialName != credentialName
	})
	encodedCreds, err := json.Marshal(results)
	if err != nil {
		return backend.WrapError("Unable to delete credential", err)
	}
	auth.JSONData = helpers.Ptr(string(encodedCreds))

	bridge.UpdateAuthForUser(auth)

	return nil
}

func (a WebAuthn) updateCredentialName(info WebAuthnUpdateCredentialInfo, bridge authschemes.AShirtAuthBridge) error {
	userAuth, err := bridge.FindUserAuthByUserID(info.UserID)
	if err != nil {
		return backend.WrapError("Unable to find user", err)
	}
	webauthRawCreds := []byte(*userAuth.JSONData)
	var creds []AShirtWebauthnCredential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return backend.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}
	matchingIndex, _ := helpers.Find(creds, func(item AShirtWebauthnCredential) bool {
		return string(item.CredentialName) == string(info.CredentialName)
	})
	if matchingIndex == -1 {
		return backend.WrapError("Could not find matching credential", err)
	}
	creds[matchingIndex].CredentialName = info.NewCredentialName
	creds[matchingIndex].CredentialCreatedDate = time.Now()

	encodedCreds, err := json.Marshal(creds)
	if err != nil {
		return backend.WrapError("Unable to encode credentials", err)
	}
	userAuth.JSONData = helpers.Ptr(string(encodedCreds))
	if err = bridge.UpdateAuthForUser(userAuth); err != nil {
		return backend.WrapError("Unable to update credential", err)
	}
	return nil
}

func (a WebAuthn) beginRegistration(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, info WebAuthnRegistrationInfo) (*protocol.CredentialCreation, error) {
	var user webauthnUser
	if info.RegistrationType == CreateCredential {
		user = makeNewWebAuthnUser(info.FirstName, info.LastName, info.Email, info.Username, info.CredentialName)
	} else if info.RegistrationType == LinkCredential {
		user = makeLinkingWebAuthnUser(info.UserID, info.Username, info.CredentialName)
	} else { // Add Credential
		user = makeAddCredentialWebAuthnUser(info.UserID, info.CredentialName, info.Username, info.ExistingCredentials)
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

func (a WebAuthn) beginLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, username string) (interface{}, error) {
	authData, err := bridge.FindUserAuth(username)
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

	webauthnUser := makeWebAuthnUser(user.FirstName, user.LastName, username, user.Email, user.ID, authData.AuthnID, creds)
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

	data.UserData.Credentials = append(data.UserData.Credentials, wrapCredential(*cred, AShirtWebauthnExtension{
		CredentialName:        data.UserData.CredentialName,
		CredentialCreatedDate: data.UserData.CredentialCreatedDate,
	}))

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
