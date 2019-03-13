// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/theparanoids/ashirt/backend/authschemes"
	"github.com/theparanoids/ashirt/backend/authschemes/localauth"
	"github.com/theparanoids/ashirt/backend/authschemes/oktaauth"
	"github.com/theparanoids/ashirt/backend/authschemes/recoveryauth"
	"github.com/theparanoids/ashirt/backend/config"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/logging"
	"github.com/theparanoids/ashirt/backend/server"
	"github.com/pkg/errors"
)

func main() {
	err := config.LoadWebConfig()
	logger := logging.SetupStdoutLogging()
	if err != nil {
		logger.Log("error", err, "msg", "Unable to start due to configuration error")
		return
	}
	for {
		err := tryRunServer(logger)
		logger.Log("error", err, "msg", "Restarting app")
		time.Sleep(3 * time.Second)
	}
}

func tryRunServer(logger logging.Logger) error {
	db, err := database.NewConnection(config.DBUri(), "./migrations")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Unable to connect to database (DB_URI=%s)", config.DBUri()))
	}

	logger.Log("msg", "checking database schema")
	if err := db.CheckSchema(); err != nil {
		return err
	}

	contentStore, err := contentstore.NewDevStore()
	if err != nil {
		return err
	}
	schemes := []authschemes.AuthScheme{
		recoveryauth.New(config.RecoveryExpiry()),
	}
	for _, svc := range config.SupportedAuthServices() {
		switch svc {
		case "ashirt":
			schemes = append(schemes, localauth.LocalAuthScheme{})
		case "okta":
			schemes = append(schemes, oktaauth.NewFromConfig(
				config.AuthConfigInstance(svc),
				func(map[string]string) bool {
					return true //everyone can join!
				}))
		}
	}

	http.Handle("/web/", http.StripPrefix("/web", server.Web(
		db, contentStore, &server.WebConfig{
			CSRFAuthKey:      []byte("DEVELOPMENT_CSRF_AUTH_KEY_SECRET"),
			SessionStoreKey:  []byte("DEVELOPMENT_SESSION_STORE_KEY_SECRET"),
			UseSecureCookies: false,
			AuthSchemes:      schemes,
			Logger:           logger,
		},
	)))
	http.Handle("/api/", server.API(
		db, contentStore, logger,
	))

	logger.Log("port", config.Port(), "msg", "Now Serving")
	return http.ListenAndServe(":"+config.Port(), nil)
}
