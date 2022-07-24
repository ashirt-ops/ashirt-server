package webauthn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/webauthn/constants"
	"github.com/theparanoids/ashirt-server/backend/config"
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
	if err != nil {
		return WebAuthn{}, err
	}

	if host, _, err = net.SplitHostPort(parsedUrl.Host); err != nil {
		host = parsedUrl.Host
	}

	web, err := auth.New(&auth.Config{
		RPDisplayName: cfg.DisplayName,
		RPID:          host,
	})
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
				return nil, fmt.Errorf("registration is closed to users")
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

	remux.Route(r, "POST", "/register/end", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		// opting to ignore registration flags here

		rawData := bridge.ReadAuthSchemeSession(r)
		fmt.Println("Raw Data: ", rawData)
		data, ok := rawData.(*preRegistrationSessionData)
		if !ok {
			return nil, errors.New("Unable to complete registration -- session not found or corrupt")
		}
		fmt.Printf("data: %+v\n", data)

		// out of bounds area
		bodyCopy, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Unable to copy body", err)
			return nil, err
		}
		bodyCopyStream := io.NopCloser(bytes.NewReader(bodyCopy))
		r.Body = io.NopCloser(bytes.NewReader(bodyCopy))

		fmt.Println("Received body: ", string(bodyCopy))

		_, err = protocol.ParseCredentialCreationResponseBody(bodyCopyStream)
		if err != nil {
			fmt.Println("Unable to parse response body")
			text, outerr := json.Marshal(err)
			fmt.Println("Reported error: ", string(text))
			fmt.Println("json parse error", outerr)
		}
		
		// ^^^^ out of bounds area

		cred, err := a.Web.FinishRegistration(&data.UserData, *data.WebAuthNSessionData, r)
		if err != nil {
			return nil, backend.WrapError("Unable to complete registration", err)
		}
		fmt.Printf("Cred: %+v\n", cred)
		data.UserData.Credentials = append(data.UserData.Credentials, *cred)

		userProfile := authschemes.UserProfile{
			FirstName: data.UserData.FirstName(),
			LastName:  data.UserData.LastName(),
			Slug:      data.UserData.Email(),
			Email:     data.UserData.Email(),
		}

		userResult, err := bridge.CreateNewUser(userProfile)

		if err != nil {
			return nil, backend.WrapError("Unable to complete registration (2)", err)
		}

		return nil, bridge.CreateNewAuthForUser(authschemes.UserAuthData{
			UserID:  userResult.UserID,
			UserKey: userProfile.Email,
		})
	}))

	remux.Route(r, "POST", "/login", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		// TODO
		return nil, nil
	}))

	remux.Route(r, "POST", "/finishlogin", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		// TODO
		return nil, nil
	}))
}

func (a WebAuthn) beginRegistration(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, info WebAuthnRegistrationInfo) (*protocol.CredentialCreation, error) {

	user := makeWebAuthnUser(info.FirstName, info.LastName, info.Email)
	credData, sessionData, err := a.Web.BeginRegistration(&user) // TODO: do we want any options?

	if err != nil {
		return nil, err
	}

	err = bridge.SetAuthSchemeSession(w, r, &preRegistrationSessionData{
		UserData:            user,
		WebAuthNSessionData: sessionData,
	})

	fmt.Println("Session save error:", err)

	return credData, err
}
