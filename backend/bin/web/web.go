package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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

func main() {
	err := config.LoadWebConfig()
	logger := logging.SetupStdoutLogging()
	if err != nil {
		logging.Fatal(logger, "Unable to start due to configuration error", "error", err, "action", "exiting")
	}

	db, err := database.NewConnection(config.DBUri(), "/migrations")
	if err != nil {
		logging.Fatal(logger, "Unable to connect to database", "error", err, "action", "exiting")
	}

	logger.Info("checking database schema")
	if err := db.CheckSchema(); err != nil {
		logging.Fatal(logger, "schema read error", "error", err)
	}

	contentStore, err := confighelpers.ChooseContentStoreType(config.AllStoreConfig())
	if errors.Is(err, backend.ErrorDeprecated) {
		logger.Warn("No content store provided")
		contentStore, err = confighelpers.DefaultS3Store()
	}
	if err != nil {
		logging.Fatal(logger, "store setup error", "error", err)
	}
	logger.Info("Using Storage", "type", contentStore.Name())

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
			logger.Error("Unable to load auth scheme. Disabling.",
				"schemeName", schemeError.name,
				"error", schemeError.err.Error())
		}
		// return fmt.Errorf("Cannot continue with auth scheme failures") // Not sure if we want to just now allow certain schemes if they fail, or outright fail to launch
	}

	if config.EmailType() != "" {
		startEmailServices(db, logger)
	} else {
		logger.Warn("No Emailer selected")
	}

	r := chi.NewRouter()

	r.Route("/web", func(r chi.Router) {
		server.Web(r,
			db, contentStore, &server.WebConfig{
				SessionStoreKey:  []byte(config.SessionStoreKey()),
				UseSecureCookies: true,
				AuthSchemes:      schemes,
				Logger:           logger,
			},
		)
	})

	r.Route("/api", func(r chi.Router) {
		server.API(r,
			db, contentStore, logger,
		)
	})

	static := config.Static()
	if static == "" {
		logging.Fatal(logger, "no static directory provided")
	} else {
		logger.Info("serving static directory", "path", static)
	}

	r.Mount("/", http.FileServer(http.Dir(static)))

	logger.Info("starting Web server", "port", config.Port())
	serveErr := http.ListenAndServe(":"+config.Port(), r)
	logging.Fatal(logger, "server shutting down", "err", serveErr)
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

func startEmailServices(db *database.Connection, logger *slog.Logger) {
	var emailServicer emailservices.EmailServicer
	emailLogger := logger.With("service", "email-sender", "type", config.EmailType)
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
		logger.Error("unsupported emailer", "type", config.EmailType)
	} else {
		emailLogger.Info("Staring emailer")
		emailWorker := workers.MakeEmailWorker(db, emailServicer, logger.With("service", "email-worker"))
		emailWorker.Start()
	}
}
