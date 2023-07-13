// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	recoveryConsts "github.com/ashirt-ops/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/ashirt-ops/ashirt-server/backend/session"
)

type WebConfig struct {
	DBConnection     *database.Connection
	AuthSchemes      []authschemes.AuthScheme
	CSRFAuthKey      []byte
	SessionStoreKey  []byte
	UseSecureCookies bool
	Logger           logging.Logger
}

func (c *WebConfig) validate() error {
	if c.Logger == nil {
		fmt.Println(`error="Logger not set" action="Using NopLogger"`)
		c.Logger = logging.NewNopLogger()
	}
	if len(c.CSRFAuthKey) < 32 {
		return errors.New("CSRFAuthKey must be 32 bytes or longer")
	}
	if len(c.SessionStoreKey) < 32 {
		return errors.New("SessionStoreKey must be 32 bytes or longer")
	}
	if !c.UseSecureCookies {
		c.Logger.Log("msg", "Config Warning: cookies not using secure flag")
	}
	return nil
}

func Web(r chi.Router, db *database.Connection, contentStore contentstore.Store, config *WebConfig) {
	if err := config.validate(); err != nil {
		panic(err)
	}
	sessionStore, err := session.NewStore(db, session.StoreOptions{
		SessionDuration:  30 * 24 * time.Hour,
		UseSecureCookies: config.UseSecureCookies,
		Key:              config.SessionStoreKey,
	})
	if err != nil {
		panic(err)
	}

	r.Handle("/metrics", promhttp.Handler())
	r.Group(func(r chi.Router) {
		r.Use(middleware.LogRequests(config.Logger))
		r.Use(csrf.Protect(config.CSRFAuthKey,
			csrf.Secure(config.UseSecureCookies),
			csrf.Path("/"),
			csrf.ErrorHandler(jsonHandler(func(r *http.Request) (interface{}, error) {
				return nil, backend.CSRFErr(csrf.FailureReason(r))
			}))))
		r.Use(middleware.InjectCSRFTokenHeader())
		r.Use(middleware.AuthenticateUserAndInjectCtx(db, sessionStore))

		supportedAuthSchemes := make([]dtos.SupportedAuthScheme, len(config.AuthSchemes))
		for i, scheme := range config.AuthSchemes {
			r.Route("/auth/"+scheme.Name(), func(r chi.Router) {
				scheme.BindRoutes(r.(chi.Router), authschemes.MakeAuthBridge(db, sessionStore, scheme.Name(), scheme.Type()))
			})
			supportedAuthSchemes[i] = dtos.SupportedAuthScheme{
				SchemeName:  scheme.FriendlyName(),
				SchemeCode:  scheme.Name(),
				SchemeFlags: scheme.Flags(),
				SchemeType:  scheme.Type(),
			}
		}
		authsWithOutRecovery := make([]dtos.SupportedAuthScheme, 0, len(supportedAuthSchemes)-1)

		// recovery is a special authentication that we kind of want to hide/separate from the other auth schemes
		// so, we filter it out here
		for _, auth := range supportedAuthSchemes {
			if auth.SchemeCode != recoveryConsts.Code {
				authsWithOutRecovery = append(authsWithOutRecovery, auth)
			}
		}

		bindSharedRoutes(r, db, contentStore)
		bindWebRoutes(r, db, contentStore, sessionStore, &authsWithOutRecovery)
	})
}

func bindWebRoutes(r chi.Router, db *database.Connection, contentStore contentstore.Store, sessionStore *session.Store, supportedAuthSchemes *[]dtos.SupportedAuthScheme) {
	route(r, "POST", "/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonHandler(func(r *http.Request) (interface{}, error) {
			err := sessionStore.Delete(w, r)
			if err != nil {
				return nil, backend.WrapError("Unable to delete session", err)
			}
			return nil, nil
		}).ServeHTTP(w, r)
	}))

	route(r, "GET", "/user", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		slug := dr.FromQuery("userSlug").AsString()

		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ReadUser(r.Context(), db, slug, supportedAuthSchemes)
	}))

	route(r, "GET", "/auths", jsonHandler(func(r *http.Request) (interface{}, error) {
		return supportedAuthSchemes, nil
	}))

	route(r, "GET", "/auths/breakdown", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListAuthDetails(r.Context(), db, supportedAuthSchemes)
	}))

	route(r, "POST", "/operations/{operation_slug}/evidence", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectFormRequest(r)
		i := services.CreateEvidenceInput{
			Description:   dr.FromBody("description").Required().AsString(),
			Content:       dr.FromFile("content"),
			ContentType:   dr.FromBody("contentType").OrDefault("image").AsString(),
			OccurredAt:    dr.FromBody("occurredAt").OrDefault(time.Now()).AsTime(),
			OperationSlug: dr.FromURL("operation_slug").AsString(),
		}
		tagIDsJSON := dr.FromBody("tagIds").OrDefault("[]").AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		if err := json.Unmarshal([]byte(tagIDsJSON), &i.TagIDs); err != nil {
			return nil, backend.BadInputErr(err, "tagIds must be a json array of ints")
		}
		return services.CreateEvidence(r.Context(), db, contentStore, i)
	}))
}
