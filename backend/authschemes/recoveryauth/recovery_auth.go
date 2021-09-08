// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package recoveryauth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/server/remux"
)

type RecoveryAuthScheme struct {
	Expiry time.Duration
}

func New(maxAge time.Duration) RecoveryAuthScheme {
	return RecoveryAuthScheme{Expiry: maxAge}
}

// Name returns the name of this authscheme
func (RecoveryAuthScheme) Name() string {
	return constants.Code
}

// FriendlyName returns "ASHIRT User Recovery"
func (RecoveryAuthScheme) FriendlyName() string {
	return constants.FriendlyName
}

// Flags returns an empty string (no supported auth flags for recovery)
func (RecoveryAuthScheme) Flags() []string {
	return []string{}
}

func (RecoveryAuthScheme) Type() string {
	return constants.Code
}

func (p RecoveryAuthScheme) BindRoutes(r *mux.Router, bridge authschemes.AShirtAuthBridge) {
	remux.Route(r, "POST", "/generate", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		userSlug := dr.FromBody("userSlug").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}

		return generateRecoveryCodeForUser(r.Context(), bridge, userSlug)
	}))

	remux.Route(r, "POST", "/generateemail", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		dr := remux.DissectJSONRequest(r)
		userEmail := dr.FromBody("userEmail").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}

		err := generateRecoveryEmail(r.Context(), bridge, userEmail)
		if err != nil {
			logging.Log(r.Context(), "msg", "Unable to generate recovery email", "error", err.Error())
		}
		return nil, nil
	}))

	remux.Route(r, "POST", "/admin/register", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		if !middleware.IsAdmin(r.Context()) {
			return nil, backend.UnauthorizedWriteErr(fmt.Errorf("Requesting user is not an admin"))
		}

		dr := remux.DissectJSONRequest(r)
		firstName := dr.FromBody("firstName").Required().AsString()
		lastName := dr.FromBody("lastName").AsString()
		profile := authschemes.UserProfile{
			FirstName: firstName,
			LastName:  lastName,
			Slug:      strings.ToLower(firstName + "." + lastName),
			Email:     dr.FromBody("email").Required().AsString(),
		}

		if dr.Error != nil {
			return nil, dr.Error
		}
		userResult, err := bridge.CreateNewUser(profile)
		if err != nil {
			return nil, backend.WrapError("Unable to create new user", err)
		}

		return generateRecoveryCodeForUser(r.Context(), bridge, userResult.RealSlug)
	}))

	remux.Route(r, "GET", "/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dr := remux.DissectJSONRequest(r)
		recoveryKey := dr.FromQuery("code").Required().AsString()

		if dr.Error != nil {
			remux.HandleError(w, r, dr.Error)
			return
		}

		userID, err := bridge.OneTimeVerification(r.Context(), recoveryKey, int64(p.Expiry/time.Minute))
		if err != nil {
			http.Redirect(w, r, "/autherror/recoveryfailed", http.StatusFound)
			return
		}
		bridge.LoginUser(w, r, userID, nil)
		http.Redirect(w, r, fmt.Sprintf("/account/authmethods"), http.StatusFound)
	}))

	remux.Route(r, "DELETE", "/expired", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		return nil, DeleteExpiredRecoveryCodes(r.Context(), bridge.GetDatabase(), int64(p.Expiry/time.Minute))
	}))

	remux.Route(r, "GET", "/metrics", remux.JSONHandler(func(r *http.Request) (interface{}, error) {
		return getRecoveryMetrics(r.Context(), bridge.GetDatabase(), int64(p.Expiry/time.Minute))
	}))

}
