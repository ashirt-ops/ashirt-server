package main

import (
	stderrors "errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ashirt-ops/ashirt-server/internal/authschemes"
	"github.com/ashirt-ops/ashirt-server/internal/authschemes/localauth"
	"github.com/ashirt-ops/ashirt-server/internal/authschemes/oidcauth"
	"github.com/ashirt-ops/ashirt-server/internal/authschemes/recoveryauth"
	"github.com/ashirt-ops/ashirt-server/internal/authschemes/webauthn"
	"github.com/ashirt-ops/ashirt-server/internal/config"
	"github.com/ashirt-ops/ashirt-server/internal/config/confighelpers"
	"github.com/ashirt-ops/ashirt-server/internal/contentstore"
	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/database/seeding"
	"github.com/ashirt-ops/ashirt-server/internal/emailservices"
	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/ashirt-ops/ashirt-server/internal/logging"
	"github.com/ashirt-ops/ashirt-server/internal/server"
	"github.com/ashirt-ops/ashirt-server/internal/workers"
	"github.com/go-chi/chi/v5"

	sq "github.com/Masterminds/squirrel"
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

	db, err := database.NewConnection(config.DBUri(), config.MigrationsPath())
	if err != nil {
		logging.Fatal(logger, "Unable to connect to database", "error", err, "action", "exiting")
	}

	logger.Info("checking database schema")
	if err := db.CheckSchema(); err != nil {
		logging.Fatal(logger, "schema read error", "error", err)
	}

	seedFiles := false
	if config.SeedDatabase() {
		if seeded, err := seeding.IsSeeded(db); !seeded && err == nil {
			logger.Info("applying db seeding")
			if err := seeding.HarryPotterSeedData.ApplyTo(db); err != nil {
				logging.Fatal(logger, "seeding error", "error", err)
			}
			seedFiles = true
		}
	}

	contentStore, err := confighelpers.ChooseContentStoreType(config.AllStoreConfig())
	if stderrors.Is(err, errors.ErrorDeprecated) {
		logger.Warn("No content store provided")
		contentStore, err = confighelpers.DefaultS3Store()
	}
	if err != nil {
		logging.Fatal(logger, "store setup error", "error", err)
	}
	logger.Info("Using Storage", "type", contentStore.Name())

	if seedFiles {
		logger.Info("Adding files to storage")
		if contentStore.Name() != "local" {
			seedEvidenceFiles(db, contentStore, logger)
		}
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
				UseSecureCookies: config.UseSecureCookies(),
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

func seedEvidenceFiles(db *database.Connection, dstStore contentstore.Store, logger *slog.Logger) {
	readStore, err := confighelpers.DefaultDevStore()
	if err != nil {
		panic("Cannot create temporary devstore for copying evidence")
	}

	type evidence struct {
		FullKey  string `db:"full_image_key"`
		ThumbKey string `db:"thumb_image_key"`
	}
	var evidenceData []evidence
	err = db.Select(&evidenceData, sq.Select(
		"full_image_key", "thumb_image_key").
		From("evidence").
		Where(sq.NotEq{"content_type": "none"}),
	)

	if err != nil {
		panic("Cannot fetch evidence")
	}

	evidenceList := map[string]bool{}
	for _, evidenceItem := range evidenceData {
		evidenceList[evidenceItem.FullKey] = true
		evidenceList[evidenceItem.ThumbKey] = true
	}

	for k := range evidenceList {
		_, foundErr := dstStore.Read(k)
		if foundErr != nil {
			logger.Info("Moving content", "key", k)
			data, _ := readStore.Read(k)
			dstStore.UploadWithName(k, data)
		}
	}
}
