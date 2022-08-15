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

	if port != "" {
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
		rawData := bridge.ReadAuthSchemeSession(r)

		data, ok := rawData.(*webAuthNSessionData)
		if !ok {
			return nil, errors.New("Unable to complete registration -- session not found or corrupt")
		}

		cred, err := a.Web.FinishRegistration(&data.UserData, *data.WebAuthNSessionData, r)
		if err != nil {
			return nil, backend.WrapError("Unable to complete registration", err)
		}

		data.UserData.Credentials = append(data.UserData.Credentials, *cred)
		encodedCreds, err := json.Marshal(data.UserData.Credentials)
		if err != nil {
			return nil, backend.WrapError("Unable to create registration", err)
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
			if _, err := a.Web.FinishLogin(user, *data.WebAuthNSessionData, r); err != nil {
				return nil, backend.BadAuthErr(err)
			}

			if err := bridge.LoginUser(w, r, user.UserIDAsI64(), nil); err != nil {
				return nil, backend.WrapError("Attempt to finish login failed", err)
			}

			return nil, nil
		}).ServeHTTP(w, r)
	}))
}

func (a WebAuthn) beginRegistration(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, info WebAuthnRegistrationInfo) (*protocol.CredentialCreation, error) {
	user := makeNewWebAuthnUser(info.FirstName, info.LastName, info.Email)

	credData, sessionData, err := a.Web.BeginRegistration(&user) // TODO: do we want any options?
	if err != nil {
		return nil, err
	}

	err = bridge.SetAuthSchemeSession(w, r, makeWebauthNSessionData(user, sessionData))

	return credData, err
}

func (a WebAuthn) beginLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, email string) (interface{}, error) {

	// todo switch these to ConnectionProxy (instead of direct db requests)
	user, err := bridge.FindUserByEmail(email, false)
	if err != nil {
		return nil, backend.WrapError("Could not validate user", err)
	}

	authData, err := bridge.FindUserAuth(email)
	if err != nil {
		return nil, backend.WrapError("Could not validate user", err)
	}
	if authData.JSONData == nil {
		return nil, backend.WrapError("User lacks webauthn credentials", err)
	}

	webauthRawCreds := []byte(*authData.JSONData)
	var creds []auth.Credential
	if err = json.Unmarshal(webauthRawCreds, &creds); err != nil {
		return nil, backend.WrapError("Unable to parse webauthn credentials", err)
	}

	webauthnUser := makeWebAuthnUser(user.FirstName, user.LastName, user.Slug, user.Email, user.ID, creds)
	options, sessionData, err := a.Web.BeginLogin(&webauthnUser)
	if err != nil {
		return nil, backend.WrapError("Unable to begin login process", err)
	}

	if err = bridge.SetAuthSchemeSession(w, r, makeWebauthNSessionData(webauthnUser, sessionData)); err != nil {
		return nil, backend.WrapError("Unable to begin login process", err)
	}

	return options, nil
}
