package webauthn

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	stderrors "errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ashirt-ops/ashirt-server/internal/authschemes"
	"github.com/ashirt-ops/ashirt-server/internal/authschemes/webauthn/constants"
	"github.com/ashirt-ops/ashirt-server/internal/config"
	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/server/middleware"
	"github.com/ashirt-ops/ashirt-server/internal/server/remux"
	"github.com/go-chi/chi/v5"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthn struct {
	RegistrationEnabled bool
	Web                 *webauthn.WebAuthn
}

func New(cfg config.AuthInstanceConfig, webConfig *config.WebConfig) (WebAuthn, error) {
	parsedUrl, err := url.Parse(webConfig.FrontendIndexURL)
	if err != nil {
		return WebAuthn{}, err
	}

	host, _, err := net.SplitHostPort(parsedUrl.Host)
	if err != nil {
		return WebAuthn{}, err
	}

	rpID := host
	if cfg.WebauthnConfig.RPID != "" {
		rpID = cfg.WebauthnConfig.RPID
	}

	rpOrigins := []string{webConfig.FrontendIndexURL}
	if len(cfg.WebauthnConfig.RPOrigins) > 0 {
		rpOrigins = append(rpOrigins, cfg.WebauthnConfig.RPOrigins...)
	}

	webauthnConfig := &webauthn.Config{
		RPDisplayName: cfg.WebauthnConfig.DisplayName,
		RPID:          rpID,
		RPOrigins:     rpOrigins,
		// the below are all optional
		Debug:                  cfg.WebauthnConfig.Debug,
		AttestationPreference:  cfg.WebauthnConfig.Conveyance(),
		AuthenticatorSelection: cfg.BuildAuthenticatorSelection(),
	}

	web, err := webauthn.New(webauthnConfig)
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

// Using DissectJSONRequest(r) on discoverable requests causes the request to be parsed twice, which causes an error
// so this function allows us to get query ars without dissecting it
func isDiscoverable(r *http.Request) bool {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		return false
	}

	return parsedURL.Query().Get("discoverable") == "true"
}

func (a WebAuthn) BindRoutes(r chi.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "POST", "/register/begin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			// validate basic registration data
			if !a.RegistrationEnabled {
				return nil, stderrors.New("registration is closed to users")
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
			return nil, errors.WrapError("Unable to validate registration data", err)
		}

		userProfile := authschemes.UserProfile{
			FirstName: data.UserData.FirstName,
			LastName:  data.UserData.LastName,
			Slug:      strings.ToLower(data.UserData.FirstName + "." + data.UserData.LastName),
			Email:     data.UserData.Email,
		}
		userResult, err := bridge.CreateNewUser(userProfile)
		if err != nil {
			return nil, errors.WrapError("Unable to create user", err)
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
				return nil, stderrors.New("Unable to complete login -- session not found or corrupt")
			}
			discoverable := isDiscoverable(r)

			var cred *webauthn.Credential
			var err error

			if discoverable {
				parsedResponse, err := protocol.ParseCredentialRequestResponse(r)
				if err != nil {
					return nil, errors.WrapError("error parsing credential", errors.BadAuthErr(err))
				}

				var webauthnUser webauthnUser
				userHandler := func(_, userHandle []byte) (user webauthn.User, err error) {
					authnID := string(userHandle)
					dbUser, err := bridge.GetUserFromAuthnID(authnID)
					if err != nil {
						return nil, errors.WebauthnLoginError(err, "Could not find user from authn ID", "No such user found")
					}
					auth, err := bridge.FindUserAuthByUserID(dbUser.ID)
					if err != nil {
						return nil, errors.DatabaseErr(err)
					}
					creds, err := a.getExistingCredentials(auth)
					if err != nil {
						return nil, err
					}
					webauthnUser = makeWebAuthnUser(dbUser.FirstName, dbUser.LastName, dbUser.Slug, dbUser.Email, dbUser.ID, userHandle, creds)
					return &webauthnUser, nil
				}
				cred, err = a.Web.ValidateDiscoverableLogin(userHandler, *data.WebAuthNSessionData, parsedResponse)
				if err != nil {
					return nil, errors.BadAuthErr(err)
				} else if cred.Authenticator.CloneWarning {
					return nil, errors.WrapError("credential appears to be cloned", errors.BadAuthErr(err))
				}

				err = bridge.SetAuthSchemeSession(w, r, makeWebauthNSessionData(webauthnUser, data.WebAuthNSessionData))
				if err != nil {
					return nil, errors.WebauthnLoginError(err, "Unable to finish login process", "Unable to set session")
				}
				rawData = bridge.ReadAuthSchemeSession(r)
				data, ok = rawData.(*webAuthNSessionData)
				if !ok {
					return nil, stderrors.New("Unable to finish login -- session not found or corrupt")
				}
			} else {
				user := &data.UserData
				cred, err = a.Web.FinishLogin(user, *data.WebAuthNSessionData, r)
				if err != nil {
					return nil, errors.BadAuthErr(err)
				} else if cred.Authenticator.CloneWarning {
					return nil, errors.WrapError("credential appears to be cloned", errors.BadAuthErr(err))
				}
			}

			updateSignCount(data, cred, bridge)

			if err := bridge.LoginUser(w, r, data.UserData.UserIDAsI64(), nil); err != nil {
				return nil, errors.WrapError("Attempt to finish login failed", err)
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
			return nil, errors.WrapError("Unable to validate registration data", err)
		}

		rawSessionData := bridge.ReadAuthSchemeSession(r)
		sessionData, _ := rawSessionData.(*webAuthNSessionData)

		return nil, bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:   byteSliceToI64(data.UserData.UserID),
			AuthnID:  sessionData.UserData.AuthnID,
			Username: data.UserData.UserName,
			JSONData: helpers.Ptr(string(encodedCreds)),
		})
	}))

	remux.Route(r, "GET", "/credentials", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		callingUserID := middleware.UserID(r.Context())
		return a.getCredentials(callingUserID, bridge)
	}))

	remux.Route(r, "DELETE", "/credential/{credentialID}", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		callingUserID := middleware.UserID(r.Context())
		dr := remux.DissectJSONRequest(r)
		credentialID := dr.FromURL("credentialID").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		credIDByteArr, _ := hex.DecodeString(credentialID)
		return nil, a.deleteCredential(callingUserID, credIDByteArr, bridge)
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
				return nil, errors.DatabaseErr(err)
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
			return nil, errors.WrapError("Unable to validate registration data", err)
		}

		userAuth, err := bridge.FindUserAuthByContext(r.Context())
		if err != nil {
			return nil, errors.WrapError("Unable to find user", err)
		}
		userAuth.JSONData = helpers.Ptr(string(encodedCreds))
		err = bridge.UpdateAuthForUser(userAuth)
		if err != nil {
			return nil, errors.WrapError("Unable to update credentials", err)
		}

		// We might want to return a full list of credentials. TODO: check if we want that
		return nil, nil
	}))
}

func (a WebAuthn) getCredentials(userID int64, bridge authschemes.AShirtAuthBridge) (*ListCredentialsOutput, error) {
	auth, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return nil, errors.WrapError("Unable to get credentials", err)
	}

	webauthRawCreds := []byte(*auth.JSONData)
	var creds []AShirtWebauthnCredential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return nil, errors.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	results := helpers.Map(creds, func(cred AShirtWebauthnCredential) CredentialEntry {
		return CredentialEntry{
			CredentialName: cred.CredentialName,
			DateCreated:    cred.CredentialCreatedDate,
			CredentialID:   hex.EncodeToString(cred.ID),
		}
	})
	output := ListCredentialsOutput{results}
	return &output, nil
}

func (a WebAuthn) deleteCredential(userID int64, credentialID []byte, bridge authschemes.AShirtAuthBridge) error {
	auth, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return errors.WrapError("Unable to find user", err)
	}

	webauthRawCreds := []byte(*auth.JSONData)
	var creds []AShirtWebauthnCredential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return errors.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	results := helpers.Filter(creds, func(cred AShirtWebauthnCredential) bool {
		return !bytes.Equal(cred.ID, credentialID)
	})
	encodedCreds, err := json.Marshal(results)
	if err != nil {
		return errors.WrapError("Unable to delete credential", err)
	}
	auth.JSONData = helpers.Ptr(string(encodedCreds))

	bridge.UpdateAuthForUser(auth)

	return nil
}

func (a WebAuthn) updateCredentialName(info WebAuthnUpdateCredentialInfo, bridge authschemes.AShirtAuthBridge) error {
	userAuth, err := bridge.FindUserAuthByUserID(info.UserID)
	if err != nil {
		return errors.WrapError("Unable to find user", err)
	}
	webauthRawCreds := []byte(*userAuth.JSONData)
	var creds []AShirtWebauthnCredential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return errors.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}
	matchingIndex, _ := helpers.Find(creds, func(item AShirtWebauthnCredential) bool {
		return string(item.CredentialName) == string(info.CredentialName)
	})
	if matchingIndex == -1 {
		return errors.WrapError("Could not find matching credential", err)
	}
	creds[matchingIndex].CredentialName = info.NewCredentialName
	creds[matchingIndex].CredentialCreatedDate = time.Now()

	encodedCreds, err := json.Marshal(creds)
	if err != nil {
		return errors.WrapError("Unable to encode credentials", err)
	}
	userAuth.JSONData = helpers.Ptr(string(encodedCreds))
	if err = bridge.UpdateAuthForUser(userAuth); err != nil {
		return errors.WrapError("Unable to update credential", err)
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
		user = makeAddCredentialWebAuthnUser(info.UserID, info.Username, info.CredentialName, info.ExistingCredentials)
	}

	idx, _ := helpers.Find(info.ExistingCredentials, func(cred AShirtWebauthnCredential) bool {
		return strings.ToLower(cred.CredentialName) == strings.ToLower(info.CredentialName)
	})

	if idx != -1 {
		return nil, errors.BadInputErr(
			stderrors.New("user trying to register with taken credential name"),
			"Credential name is already taken",
		)
	}

	discoverable := isDiscoverable(r)

	credExcludeList := make([]protocol.CredentialDescriptor, len(user.Credentials))
	for i, cred := range user.Credentials {
		credExcludeList[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
	}

	var selection protocol.AuthenticatorSelection

	if discoverable {
		selection = protocol.AuthenticatorSelection{
			ResidentKey: protocol.ResidentKeyRequirementRequired,
		}
	}

	registrationOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = credExcludeList
		credCreationOpts.AuthenticatorSelection = selection
	}

	credOptions, sessionData, err := a.Web.BeginRegistration(&user, webauthn.WithAuthenticatorSelection(selection), registrationOptions)
	if err != nil {
		return nil, err
	}

	err = bridge.SetAuthSchemeSession(w, r, makeWebauthNSessionData(user, sessionData))

	return credOptions, err
}

func (a WebAuthn) beginLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, username string) (interface{}, error) {
	discoverable := isDiscoverable(r)

	var data interface{}
	var options *protocol.CredentialAssertion
	var sessionData *webauthn.SessionData
	var err error

	if discoverable {
		var opts = []webauthn.LoginOption{
			webauthn.WithUserVerification(protocol.VerificationPreferred),
		}
		options, sessionData, err = a.Web.BeginDiscoverableLogin(opts...)

		if err != nil {
			return nil, errors.WebauthnLoginError(err, "Unable to find login credentials", "Unable to find login credentials")
		}
		data = makeDiscoverableWebauthNSessionData(sessionData)
	} else {
		authData, err := bridge.FindUserAuth(username)
		if err != nil {
			return nil, errors.WebauthnLoginError(err, "Could not validate user", "No such auth")
		}
		if authData.JSONData == nil {
			return nil, errors.WebauthnLoginError(err, "User lacks webauthn credentials")
		}

		user, err := bridge.GetUserFromID(authData.UserID)
		if err != nil {
			return nil, errors.WebauthnLoginError(err, "Could not validate user", "No such user")
		}

		creds, err := a.getExistingCredentials(authData)
		if err != nil {
			return nil, errors.WebauthnLoginError(err, "Unable to parse webauthn credentials")
		}

		webauthnUser := makeWebAuthnUser(user.FirstName, user.LastName, username, user.Email, user.ID, authData.AuthnID, creds)
		options, sessionData, err = a.Web.BeginLogin(&webauthnUser)
		if err != nil {
			return nil, errors.WebauthnLoginError(err, "Unable to begin login process")
		}
		data = makeWebauthNSessionData(webauthnUser, sessionData)
	}

	err = bridge.SetAuthSchemeSession(w, r, data)
	if err != nil {
		return nil, errors.WebauthnLoginError(err, "Unable to begin login process", "Unable to set session")
	}

	return options, nil
}

func (a WebAuthn) getExistingCredentials(authData authschemes.UserAuthData) ([]AShirtWebauthnCredential, error) {
	webauthRawCreds := []byte(*authData.JSONData)
	var creds []AShirtWebauthnCredential
	if err := json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return nil, errors.WebauthnLoginError(err, "Unable to parse webauthn credentials")
	}

	return creds, nil
}

func (a WebAuthn) validateRegistrationComplete(r *http.Request, bridge authschemes.AShirtAuthBridge) (*webAuthNSessionData, []byte, error) {
	rawData := bridge.ReadAuthSchemeSession(r)
	data, ok := rawData.(*webAuthNSessionData)
	if !ok {
		return nil, nil, stderrors.New("Unable to complete registration -- session not found or corrupt")
	}

	cred, err := a.Web.FinishRegistration(&data.UserData, *data.WebAuthNSessionData, r)
	if err != nil {
		return nil, nil, errors.WrapError("Unable to complete registration", err)
	}

	data.UserData.Credentials = append(data.UserData.Credentials, wrapCredential(*cred, AShirtWebauthnExtension{
		CredentialName:        data.UserData.CredentialName,
		CredentialCreatedDate: data.UserData.CredentialCreatedDate,
	}))

	encodedCreds, err := json.Marshal(data.UserData.Credentials)
	if err != nil {
		return nil, nil, errors.WrapError("Unable to create registration", err)
	}
	return data, encodedCreds, nil
}

func updateSignCount(data *webAuthNSessionData, loginCred *webauthn.Credential, bridge authschemes.AShirtAuthBridge) error {
	userID := data.UserData.UserIDAsI64()

	userAuth, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return errors.WrapError("Unable to find user", err)
	}
	matchingIndex, _ := helpers.Find(data.UserData.Credentials, func(item AShirtWebauthnCredential) bool {
		return string(item.ID) == string(loginCred.ID)
	})
	if matchingIndex == -1 {
		return errors.WrapError("Could not find matching credential", err)
	}
	data.UserData.Credentials[matchingIndex].Authenticator.SignCount = loginCred.Authenticator.SignCount
	encodedCreds, err := json.Marshal(data.UserData.Credentials)
	if err != nil {
		return errors.WrapError("Unable to encode credentials", err)
	}
	userAuth.JSONData = helpers.Ptr(string(encodedCreds))
	if err = bridge.UpdateAuthForUser(userAuth); err != nil {
		return errors.WrapError("Unable to update credential", err)
	}
	return nil
}
