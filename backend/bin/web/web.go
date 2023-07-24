// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/localauth"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/oidcauth"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/recoveryauth"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/webauthn"
	"github.com/ashirt-ops/ashirt-server/backend/config"
	"github.com/ashirt-ops/ashirt-server/backend/config/confighelpers"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/emailservices"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/server"
	"github.com/ashirt-ops/ashirt-server/backend/workers"
	"github.com/go-chi/chi/v5"
)

type SchemeError struct {
	name string
	err  error
}

var sessionManager *scs.SessionManager

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

	contentStore, err := confighelpers.ChooseContentStoreType(config.AllStoreConfig())
	if errors.Is(err, backend.ErrorDeprecated) {
		logger.Log("msg", "No content store provided")
		contentStore, err = confighelpers.DefaultS3Store()
	}
	if err != nil {
		logging.Fatal(logger, "msg", "store setup error", "error", err)
	}
	logger.Log("msg", "Using Storage", "type", contentStore.Name())

	schemes := []authschemes.AuthScheme{
		recoveryauth.New(config.RecoveryExpiry()),
	}
	authSchemeNames := config.SupportedAuthServices()
	schemeErrors := make([]SchemeError, 0, len(authSchemeNames))

	for _, svc := range authSchemeNames {
		scheme, err := handleAuthType(config.AuthConfigInstance(svc))
		if err != nil {
			schemeErrors = append(schemeErrors, SchemeError{svc, err})
		} else {
			schemes = append(schemes, scheme)
		}
	}

	if len(schemeErrors) > 0 {
		for _, schemeError := range schemeErrors {
			logger.Log("msg", "Unable to load auth scheme. Disabling.",
				"schemeName", schemeError.name,
				"error", schemeError.err.Error())
		}
		// return fmt.Errorf("Cannot continue with auth scheme failures") // Not sure if we want to just now allow certain schemes if they fail, or outright fail to launch
	}

	if config.EmailType() != "" {
		startEmailServices(db, logger)
	} else {
		logger.Log("msg", "No Emailer selected")
	}

	r := chi.NewRouter()

	sessionManager = scs.New()

	r.Route("/web", func(r chi.Router) {
		server.Web(r, sessionManager,
			db, contentStore, &server.WebConfig{
				CSRFAuthKey:      []byte(config.CSRFAuthKey()),
				SessionStoreKey:  config.SessionStoreKey(),
				UseSecureCookies: true,
				AuthSchemes:      schemes,
				Logger:           logger,
			},
		)
	})

	logger.Log("msg", "starting Web server", "port", config.Port())
	serveErr := http.ListenAndServe(":"+config.Port(), sessionManager.LoadAndSave(r))
	logging.Fatal(logger, "msg", "server shutting down", "err", serveErr)
}

func handleAuthType(cfg config.AuthInstanceConfig) (authschemes.AuthScheme, error) {
	appConfig := config.AllAppConfig()
	if cfg.Type == "oidc" {
		authScheme, err := oidcauth.New(cfg, &appConfig)
		return authScheme, err
	}
	if cfg.Name == "ashirt" {
		authScheme := localauth.LocalAuthScheme{
			RegistrationEnabled: cfg.RegistrationEnabled,
		}
		return authScheme, nil
	}
	if cfg.Name == "webauthn" {
		return webauthn.New(cfg, &appConfig)
	}

	return nil, fmt.Errorf("unknown auth type: %v", cfg.Type)
}

func startEmailServices(db *database.Connection, logger logging.Logger) {
	var emailServicer emailservices.EmailServicer
	emailLogger := logging.With(logger, "service", "email-sender", "type", config.EmailType)
	switch config.EmailType() {
	case string(emailservices.StdOutEmailer):
		mailer := emailservices.MakeWriterMailer(os.Stdout, emailLogger)
		emailServicer = &mailer
	case string(emailservices.MemoryEmailer):
		mailer := emailservices.MakeMemoryMailer(emailLogger)
		emailServicer = &mailer
	case string(emailservices.SMTPEmailer):
		mailer := emailservices.MakeSMTPMailer(emailLogger)
		emailServicer = &mailer
	}

	if emailServicer == nil {
		logger.Log("msg", "unsupported emailer", "type", config.EmailType)
	} else {
		emailLogger.Log("msg", "Staring emailer")
		emailWorker := workers.MakeEmailWorker(db, emailServicer, logging.With(logger, "service", "email-worker"))
		emailWorker.Start()
	}
}
