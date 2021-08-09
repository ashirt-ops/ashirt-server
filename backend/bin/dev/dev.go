// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/authschemes/localauth"
	"github.com/theparanoids/ashirt-server/backend/authschemes/oidcauth"
	"github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/database/seeding"
	"github.com/theparanoids/ashirt-server/backend/emailservices"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/server"
	"github.com/theparanoids/ashirt-server/backend/workers"
)

type SchemeError struct {
	name string
	err  error
}

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
		return fmt.Errorf("Unable to connect to database (DB_URI=%s) : %w", config.DBUri(), err)
	}

	logger.Log("msg", "checking database schema")
	if err := db.CheckSchema(); err != nil {
		return err
	}

	if seeded, err := seeding.IsSeeded(db); !seeded && err == nil {
		logger.Log("msg", "applying db seeding")
		err := seeding.HarryPotterSeedData.ApplyTo(db)
		if err != nil {
			return err
		}
	}

	contentStore, err := contentstore.NewDevStore()
	if err != nil {
		return err
	}
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
