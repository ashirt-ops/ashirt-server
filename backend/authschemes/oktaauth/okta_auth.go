// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package oktaauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/okta/okta-jwt-verifier-golang/utils"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/server/remux"
)

// OktaAuth provides okta-based authentication via the OAuth2.0 Authorization flow.
type OktaAuth struct {
	clientID                      string
	clientSecret                  string
	issuer                        string
	absoluteBackendPath           string
	authSuccessRedirectPath       string
	authFailureRedirectPathPrefix string
	profileToShortnameField       string
	canAccessService              func(map[string]string) bool
}

// functions for mocking
type tokenVerifier func(token, nonce string) (*verifier.Jwt, error)
type codeExchanger func(code string, r *http.Request) Exchange
type profileGatherer func(r *http.Request, accessToken string) map[string]string

// New generates a new OktaAuth instance. OktaAuth requires a fair bit of configuration.
// These are the fields and their meanings:
//
// clientID: The Okta-generated client id for an application
//
// clientSecret:The Okta-generated client secret for an application
//
// issuer: The okta protocol, domain, and root path for okta verification.
//
// backendAbsolutePath: The protocol, domain, and root path for the backend (e.g. "http://localhost:3000/web")
//
// successRedirectPath: The absolute path on where to redirect the user when auth is successful
//
// failureRedirectPath: The absolute path on where to redirect the user when auth fails
//
// canAccessService: A function that evaluates an okta profile (map[string]string) to
// determine if a user has access to this application. If the user should have access to this service,
// then return true. Otherwise, return false.
func New(clientID, clientSecret, issuer, backendPath, successRedirectPath, failureRedirectPathPrefix, profileToShortnameField string, canAccessService func(map[string]string) bool) OktaAuth {
	return OktaAuth{
		clientID:                      clientID,
		clientSecret:                  clientSecret,
		issuer:                        issuer,
		absoluteBackendPath:           backendPath,
		authSuccessRedirectPath:       successRedirectPath,
		authFailureRedirectPathPrefix: failureRedirectPathPrefix,
		canAccessService:              canAccessService,
		profileToShortnameField:       profileToShortnameField,
	}
}

func NewFromConfig(cfg config.AuthInstanceConfig, canAccessService func(map[string]string) bool) OktaAuth {
	return New(cfg.ClientID, cfg.ClientSecret, cfg.Issuer, cfg.BackendURL, cfg.SuccessRedirectURL, cfg.FailureRedirectURLPrefix, cfg.ProfileToShortnameField, canAccessService)
}

// Name returns back "okta"
func (OktaAuth) Name() string {
	return "okta"
}

// FriendlyName returns "Okta OIDC"
func (OktaAuth) FriendlyName() string {
	return "Okta OIDC"
}

func (okta OktaAuth) authSuccess(w http.ResponseWriter, r *http.Request, linking bool) (interface{}, error) {
	if linking {
		return authDone(w, r, "/account/authmethods", nil)
	}
	return authDone(w, r, okta.authSuccessRedirectPath, nil)
}

func (okta OktaAuth) authFailure(w http.ResponseWriter, r *http.Request, err error, errorPath string) (interface{}, error) {
	return authDone(w, r, okta.authFailureRedirectPathPrefix+errorPath, err)
}

func authDone(w http.ResponseWriter, r *http.Request, frontendPath string, err error) (interface{}, error) {
	http.Redirect(w, r, frontendPath, http.StatusFound)
	return nil, err
}

// callbackURI provides a consistent url for callbacks
func (okta OktaAuth) callbackURI() string {
	return fmt.Sprintf("%v/auth/%v/callback", okta.absoluteBackendPath, okta.Name())
}

// makeUserProfile constructs a basic ashirt user profile from an okta v1 profile
// Note: this expects that the following values are present on an okta profile: given_name, family_name
// (plus whatever is expected from oktaProfileSlugField)
func (okta OktaAuth) makeUserProfile(profile map[string]string) authschemes.UserProfile {
	return authschemes.UserProfile{
		FirstName: profile["given_name"],
		LastName:  profile["family_name"],
		Slug:      profile[okta.profileToShortnameField],
		Email:     profile["email"],
	}
}

func (okta OktaAuth) redirectLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, mode string) {
	nonce, _ := utils.GenerateNonce()
	stateChallenge := csrf.Token(r)
	bridge.SetAuthSchemeSession(w, r, &preLoginAuthSession{
		Nonce:              nonce,
		StateChallengeCSRF: stateChallenge,
		OktaMode:           mode,
	})

	q := r.URL.Query()
	q.Add("client_id", okta.clientID)
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email") // TODO: scope may need to change with okta.profileToShortnameField
	q.Add("redirect_uri", okta.callbackURI())
	q.Add("state", stateChallenge)
	q.Add("nonce", nonce)

	http.Redirect(w, r, fmt.Sprintf("%v/v1/authorize?%v", okta.issuer, q.Encode()), http.StatusFound)
}

const (
	modeLogin = "login"
	modeLink  = "link"
)

// BindRoutes implements two routes to complete the okta login/ashirt registration process.
// /login kicks off the process, redirecting the user to okta to login. Once successful,
// okta will contact /callback to complete the process. In addition to normal auth verification,
// /callback also checks that a user is allowed to access this service (via the canAccessService function
// provided via oktaauth.New) and will generate a new ashirt user if that user doesn't already exist.
func (okta OktaAuth) BindRoutes(r *mux.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "GET", "/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		okta.redirectLogin(w, r, bridge, modeLogin)
	}))

	remux.Route(r, "GET", "/link", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		okta.redirectLogin(w, r, bridge, modeLink)
	}))

	remux.Route(r, "GET", "/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			return okta.handleOktaCallback(w, r, bridge, okta.exchangeCode, okta.verifyToken, okta.getProfileData)
		}).ServeHTTP(w, r)
	}))
}

func (okta OktaAuth) handleOktaCallback(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge,
	exchangeCode codeExchanger, verifyToken tokenVerifier, getUserProfile profileGatherer, // for mocks
) (interface{}, error) {
	oktaCode := r.URL.Query().Get("code")

	sess, ok := bridge.ReadAuthSchemeSession(r).(*preLoginAuthSession)
	if !ok {
		return okta.authFailure(w, r, backend.BadAuthErr(errors.New("Callback called without preloginauth session")), "/autherror/noaccess")
	}

	linkingAccount := sess.OktaMode == modeLink

	if r.URL.Query().Get("state") != sess.StateChallengeCSRF || oktaCode == "" {
		return okta.authFailure(w, r, backend.BadAuthErr(errors.New("Authentication challenge failed")), "/autherror/noverify")
	}

	exchange := exchangeCode(oktaCode, r)
	if exchange.WrappedError != nil {
		return okta.authFailure(w, r, exchange.WrappedError, "/autherror/noverify")
	}
	if exchange.Error != "" || exchange.ErrorDescription != "" {
		return okta.authFailure(w, r, fmt.Errorf("%v : %v", exchange.Error, exchange.ErrorDescription), "/autherror/noverify")
	}

	_, verificationError := verifyToken(exchange.IDToken, sess.Nonce)

	if verificationError != nil {
		return okta.authFailure(w, r, backend.BadAuthErr(errors.New("Authentication token verification failed")), "/autherror/noverify")
	}

	profile := getUserProfile(r, exchange.AccessToken)
	if !okta.canAccessService(profile) {
		return okta.authFailure(w, r, backend.BadAuthErr(errors.New("User is not permitted access")), "/autherror/noaccess")
	}

	shortName, ok := profile[okta.profileToShortnameField]
	if !ok || shortName == "" {
		return okta.authFailure(w, r, backend.BadAuthErr(errors.New("Shortname is empty, check that profileToShortNameField is correct")), "/autherror/noaccess")
	}

	authData, err := bridge.FindUserAuth(shortName)
	if err != nil { //an error here implies that a user doesn't yet exist
		var userID int64
		if linkingAccount {
			userID = middleware.UserID(r.Context())
		} else {
			userResult, err := bridge.CreateNewUser(okta.makeUserProfile(profile))
			if err != nil {
				return okta.authFailure(w, r, err, "/autherror/incomplete")
			}
			userID = userResult.UserID
		}

		authData = authschemes.UserAuthData{
			UserID:  userID,
			UserKey: shortName,
		}
		err = bridge.CreateNewAuthForUser(authData)
		if err != nil {
			return okta.authFailure(w, r, err, "/autherror/incomplete")
		}
	}
	if linkingAccount {
		return okta.authSuccess(w, r, linkingAccount)
	}

	err = bridge.LoginUser(w, r, authData.UserID, &authSession{
		IdToken:     exchange.IDToken,
		AccessToken: exchange.AccessToken,
	})
	if err != nil {
		if backend.IsErrorAccountDisabled(err) {
			return okta.authFailure(w, r, err, "/autherror/disabled")
		}
		return okta.authFailure(w, r, err, "/autherror/incomplete")
	}
	return okta.authSuccess(w, r, linkingAccount)
}

// the below has been adapted from https://github.com/okta/samples-golang/tree/develop/okta-hosted-login

// exchangeCode exchanges a okta-provided code for an okta provided token (which can actually do stuff)
func (okta OktaAuth) exchangeCode(code string, r *http.Request) Exchange {
	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("redirect_uri", okta.callbackURI())

	url := okta.issuer + "/v1/token?" + q.Encode()

	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		return Exchange{WrappedError: err}
	}
	h := req.Header
	req.SetBasicAuth(okta.clientID, okta.clientSecret)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Exchange{WrappedError: err}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Exchange{WrappedError: err}
	}

	defer resp.Body.Close()
	var exchange Exchange
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return Exchange{WrappedError: err}
	}
	if (exchange == Exchange{}) {
		return Exchange{WrappedError: errors.New("Unexpected okta response")}
	}

	return exchange
}

// verifyToken wraps okta-jwt-verifier-golang's JwtVerifier.VerifyIdToken function. In addition to
// standard claims, also verifies the aud matches the provided clientID, and nonce value (established in /login)
func (okta OktaAuth) verifyToken(t, nonce string) (*verifier.Jwt, error) {
	jv := verifier.JwtVerifier{
		Issuer: okta.issuer,
		ClaimsToValidate: map[string]string{
			"nonce": nonce,
			"aud":   okta.clientID,
		},
	}

	return jv.New().VerifyIdToken(t)
}

// getProfileData retrives an okta v1 profile, given an access token. If an empty access token
// is provided, returns a zero len map[string]string instead
func (okta OktaAuth) getProfileData(r *http.Request, accessToken string) map[string]string {
	profile := make(map[string]string)

	if accessToken == "" {
		logging.Log(r.Context(), "msg", "Access token not set")
		return profile
	}

	reqURL := okta.issuer + "/v1/userinfo"

	req, _ := http.NewRequest("GET", reqURL, http.NoBody)
	h := req.Header
	h.Add("Authorization", "Bearer "+accessToken)
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(body, &profile)

	return profile
}
