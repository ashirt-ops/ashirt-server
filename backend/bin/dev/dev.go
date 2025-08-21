package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/localauth"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/oidcauth"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/recoveryauth"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/webauthn"
	"github.com/ashirt-ops/ashirt-server/backend/config"
	"github.com/ashirt-ops/ashirt-server/backend/config/confighelpers"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/database/seeding"
	"github.com/ashirt-ops/ashirt-server/backend/emailservices"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/server"
	"github.com/ashirt-ops/ashirt-server/backend/workers"
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
		logger.Error("Unable to start due to configuration error", "error", err)
		return
	}
	for {
		err := tryRunServer(logger)
		logger.Error("Restarting app", "error", err)
		time.Sleep(3 * time.Second)
	}
}

func tryRunServer(logger *slog.Logger) error {
	db, err := database.NewConnection(config.DBUri(), "./migrations")
	if err != nil {
		return fmt.Errorf("Unable to connect to database (DB_URI=%s) : %w", config.DBUri(), err)
	}

	logger.Info("checking database schema")
	if err := db.CheckSchema(); err != nil {
		return err
	}

	seedFiles := false
	if seeded, err := seeding.IsSeeded(db); !seeded && err == nil {
		logger.Info("applying db seeding")
		err := seeding.HarryPotterSeedData.ApplyTo(db)
		if err != nil {
			return err
		}
		seedFiles = true
	}

	contentStore, err := confighelpers.ChooseContentStoreType(config.AllStoreConfig())
	if errors.Is(err, backend.ErrorDeprecated) {
		logger.Info("No content store provided")
		contentStore, err = confighelpers.DefaultDevStore()
	}
	if err != nil {
		return err
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
		logger.Info("No Emailer selected")
	}

	r := chi.NewRouter()

	r.Route("/web", func(r chi.Router) {
		server.Web(r,
			db, contentStore, &server.WebConfig{
				SessionStoreKey:  []byte("DEVELOPMENT_SESSION_STORE_KEY_SECRET"),
				UseSecureCookies: false,
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

	logger.Info("Now Serving", "port", config.Port())
	return http.ListenAndServe(":"+config.Port(), r)
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
