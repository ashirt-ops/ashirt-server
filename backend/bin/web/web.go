// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"net/http"

	"github.com/theparanoids/ashirt/backend/authschemes"
	"github.com/theparanoids/ashirt/backend/authschemes/localauth"
	"github.com/theparanoids/ashirt/backend/authschemes/oktaauth"
	"github.com/theparanoids/ashirt/backend/authschemes/recoveryauth"
	"github.com/theparanoids/ashirt/backend/config"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/logging"
	"github.com/theparanoids/ashirt/backend/server"
)

func main() {
	err := config.LoadWebConfig()
	logger := logging.SetupStdoutLogging()
	if err != nil {
		logging.Fatal(logger, "error", err, "msg", "Unable to start due to configuration error", "action", "exiting")
	}

	db, err := database.NewConnection(config.DBUri(), "/migrations")
	if err != nil {
		logging.Fatal(logger, "error", err, "msg", "Unable to connect to database", "action", "exiting")
	}

	logger.Log("msg", "checking database schema")
	if err := db.CheckSchema(); err != nil {
		logging.Fatal(logger, "msg", "schema read error", "error", err)
	}

	contentStore, err := contentstore.NewS3Store(config.ImageStoreBucketName(), config.ImageStoreRegion())
	if err != nil {
		logging.Fatal(logger, "msg", "image store setup error", "error", err)
	}

	archiveStore, err := contentstore.NewS3Store(config.ArchiveStoreBucketName(), config.ArchiveStoreRegion())
	if err != nil {
		logging.Fatal(logger, "msg", "archive store setup error", "error", err)
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
					return true
				}))
		}
	}

	http.Handle("/web/", http.StripPrefix("/web", server.Web(
		db, contentStore, archiveStore, &server.WebConfig{
			CSRFAuthKey:      []byte(config.CSRFAuthKey()),
			SessionStoreKey:  []byte(config.SessionStoreKey()),
			UseSecureCookies: true,
			AuthSchemes:      schemes,
			Logger:           logger,
		},
	)))

	logger.Log("msg", "starting Web server", "port", config.Port())
	serveErr := http.ListenAndServe(":"+config.Port(), nil)
	logging.Fatal(logger, "msg", "server shutting down", "err", serveErr)
}
