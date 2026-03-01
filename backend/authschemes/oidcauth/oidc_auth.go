package oidcauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	"github.com/ashirt-ops/ashirt-server/backend/config"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/ashirt-ops/ashirt-server/backend/server/remux"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

type OIDCAuth struct {
	name                          string
	friendlyName                  string
	provider                      *oidc.Provider
	oauthConfig                   oauth2.Config
	verifier                      *oidc.IDTokenVerifier
	profileSlugField              string
	profileFirstNameField         string
	profileLastNameField          string
	profileEmailField             string
	registrationEnabled           bool
	authSuccessRedirectPath       string
	authFailureRedirectPathPrefix string
}

type loginMode = string

const (
	modeLogin loginMode = "login"
	modeLink  loginMode = "link"
)

func New(cfg config.AuthInstanceConfig, webConfig *config.WebConfig) (OIDCAuth, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.ProviderURL)
	if err != nil {
		return OIDCAuth{}, err
	}

	backendURL := webConfig.BackendURL
	if cfg.BackendURL != "" {
		backendURL = cfg.BackendURL
	}

	successRedirectURL := webConfig.SuccessRedirectURL
	if cfg.SuccessRedirectURL != "" {
		successRedirectURL = cfg.SuccessRedirectURL
	}

	failureRedirectURLPrefix := cfg.FailureRedirectURLPrefix
	if cfg.FailureRedirectURLPrefix != "" {
		failureRedirectURLPrefix = cfg.FailureRedirectURLPrefix
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  callbackURI(backendURL, cfg.Name),
		Scopes:       append([]string{oidc.ScopeOpenID, "profile"}, cfg.Scopes),
		Endpoint:     provider.Endpoint(), // Discovery returns the OAuth2 endpoints.
	}

	return OIDCAuth{
		name:                          cfg.Name,
		friendlyName:                  cfg.FriendlyName,
		oauthConfig:                   oauth2Config,
		provider:                      provider,
		verifier:                      provider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID}),
		profileSlugField:              cfg.ProfileSlugField,
		profileFirstNameField:         cfg.ProfileFirstNameField,
		profileLastNameField:          cfg.ProfileLastNameField,
		profileEmailField:             cfg.ProfileEmailField,
		registrationEnabled:           cfg.RegistrationEnabled,
		authSuccessRedirectPath:       successRedirectURL,
		authFailureRedirectPathPrefix: failureRedirectURLPrefix,
	}, nil
}

func (o OIDCAuth) Name() string {
	return o.name
}

func (o OIDCAuth) FriendlyName() string {
	return o.friendlyName
}

func (OIDCAuth) Type() string {
	return "oidc"
}

// Flags returns an empty string (no supported auth flags for generic OIDC)
func (OIDCAuth) Flags() []string {
	return []string{}
}

func (o OIDCAuth) BindRoutes(r chi.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "GET", "/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.redirectLogin(w, r, bridge, modeLogin)
	}))

	remux.Route(r, "GET", "/link", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.redirectLogin(w, r, bridge, modeLink)
	}))

	remux.Route(r, "GET", "/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remux.JSONHandler(func(r *http.Request) (interface{}, error) {
			return o.handleCallback(w, r, bridge)
		}).ServeHTTP(w, r)
	}))
}

func (o OIDCAuth) redirectLogin(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge, mode string) {
	nonce, _ := authschemes.GenerateNonce()
	stateRaw := make([]byte, 16)
	io.ReadFull(rand.Reader, stateRaw)
	state := base64.RawURLEncoding.EncodeToString(stateRaw)
	bridge.SetAuthSchemeSession(w, r, &preLoginAuthSession{
		Nonce:              nonce,
		StateChallengeCSRF: state,
		LoginMode:          mode,
		OIDCService:        o.Name(),
	})
	http.Redirect(w, r, o.oauthConfig.AuthCodeURL(state), http.StatusFound)
}

func (o OIDCAuth) handleCallback(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge) (interface{}, error) {
	authName := "OIDC (" + o.friendlyName + ")"

	sess, ok := bridge.ReadAuthSchemeSession(r).(*preLoginAuthSession)
	if !ok {
		return o.authFailure(w, r, backend.BadAuthErr(errors.New(authName+" callback called without preloginauth session")), "/autherror/noaccess")
	}

	oidcExchangeCode := r.URL.Query().Get("code")
	linkingAccount := sess.LoginMode == modeLink

	if r.URL.Query().Get("state") != sess.StateChallengeCSRF || oidcExchangeCode == "" {
		return o.authFailure(w, r, backend.BadAuthErr(errors.New(authName+" authentication challenge failed")), "/autherror/noverify")
	}

	oauth2Token, err := o.oauthConfig.Exchange(r.Context(), r.URL.Query().Get("code"))

	if err != nil {
		return o.authFailure(w, r, backend.BadAuthErr(err), "/autherror/noverify")
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return o.authFailure(w, r, backend.BadAuthErr(errors.New(authName+" no id_token field in oauth2 token")), "/autherror/noverify")
	}
	idToken, err := o.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		return o.authFailure(w, r, backend.BadAuthErr(fmt.Errorf("%s authentication token verification failed: %w", authName, err)), "/autherror/noverify")
	}

	tokenSource := o.oauthConfig.TokenSource(r.Context(), oauth2Token)
	profile, err := o.provider.UserInfo(r.Context(), tokenSource)
	if err != nil {
		return o.authFailure(w, r, backend.BadAuthErr(errors.New(authName+" user is not permitted access")), "/autherror/noaccess")
	}

	profileClaims := make(map[string]interface{})
	if err = profile.Claims(&profileClaims); err != nil {
		return o.authFailure(w, r, backend.BadAuthErr(errors.New(authName+" unable to parse profile claims")), "/autherror/noaccess")
	}

	userProfile, err := o.makeUserProfile(profileClaims)
	if err != nil {
		return o.authFailure(w, r, backend.WrapError("Unable to read claim data", err), "/autherror/incomplete")
	}

	authData, err := bridge.FindUserAuth(userProfile.Slug)
	if err != nil { //an error here implies that a user doesn't yet exist
		var userID int64
		if linkingAccount {
			userID = middleware.UserID(r.Context())
		} else {
			if !o.registrationEnabled {
				return o.authFailure(w, r, backend.WrapError("Registration is disabled", err), "/autherror/registrationdisabled")
			}

			userResult, err := bridge.CreateNewUser(r.Context(), *userProfile)
			if err != nil {
				return o.authFailure(w, r, backend.WrapError("Create new "+authName+" user failed ["+userProfile.Slug+"]", err), "/autherror/incomplete")
			}
			userID = userResult.UserID
		}

		authData = authschemes.UserAuthData{
			UserID:   userID,
			Username: userProfile.Slug,
		}
		err = bridge.CreateNewAuthForUser(authData)
		if err != nil {
			return o.authFailure(w, r, backend.WrapError("Unable to create auth scheme for new "+authName+" user ["+authData.Username+"]", err), "/autherror/incomplete")
		}
	}
	if linkingAccount {
		return o.authSuccess(w, r, linkingAccount)
	}

	err = bridge.LoginUser(w, r, authData.UserID, &authSession{
		IdToken:     rawIDToken,
		AccessToken: idToken.AccessTokenHash,
	})
	if err != nil {
		if backend.IsErrorAccountDisabled(err) {
			return o.authFailure(w, r, backend.WrapError("Unable to log in "+authName+" user ["+authData.Username+"]", err), "/autherror/disabled")
		}
		return o.authFailure(w, r, backend.WrapError("Unable to log in "+authName+" user ["+authData.Username+"]", err), "/autherror/incomplete")
	}
	return o.authSuccess(w, r, linkingAccount)
}

func (o OIDCAuth) makeUserProfile(claims map[string]interface{}) (*authschemes.UserProfile, error) {
	pickValue := func(preferred, alternate string) string {
		if preferred != "" {
			return preferred
		}
		return alternate
	}
	firstNameField := pickValue(o.profileFirstNameField, "given_name")
	lastNameField := pickValue(o.profileLastNameField, "family_name")
	emailField := pickValue(o.profileEmailField, "email")
	slugField := pickValue(o.profileSlugField, "email")

	firstName, firstNameOk := claims[firstNameField].(string)
	lastName, lastNameOk := claims[lastNameField].(string)
	email, emailOk := claims[emailField].(string)
	slug, slugOk := claims[slugField].(string)

	if !all(firstNameOk, lastNameOk, emailOk, slugOk) {
		return nil, fmt.Errorf("unable to parse necessary profile fields")
	}

	userProfile := authschemes.UserProfile{
		FirstName: firstName,
		LastName:  lastName,
		Slug:      slug,
		Email:     email,
	}

	return &userProfile, nil
}

func all(fields ...bool) bool {
	isTrue := true
	for _, v := range fields {
		isTrue = isTrue && v
	}
	return isTrue
}

func (o OIDCAuth) authSuccess(w http.ResponseWriter, r *http.Request, linking bool) (interface{}, error) {
	if linking {
		return authDone(w, r, "/account/authmethods", nil)
	}
	return authDone(w, r, o.authSuccessRedirectPath, nil)
}

func (o OIDCAuth) authFailure(w http.ResponseWriter, r *http.Request, err error, errorPath string) (interface{}, error) {
	return authDone(w, r, o.authFailureRedirectPathPrefix+errorPath, err)
}

func authDone(w http.ResponseWriter, r *http.Request, frontendPath string, err error) (interface{}, error) {
	http.Redirect(w, r, frontendPath, http.StatusFound)
	return nil, err
}

func callbackURI(backendPath, name string) string {
	return fmt.Sprintf("%v/auth/%v/callback", backendPath, name)
}
