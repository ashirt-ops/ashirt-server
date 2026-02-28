package main

import (
	"errors"
	"net/http"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/config"
	"github.com/ashirt-ops/ashirt-server/backend/config/confighelpers"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/server"
	"github.com/jrozner/weby"
	webyMiddleware "github.com/jrozner/weby/middleware"
)

func main() {
	err := config.LoadAPIConfig()
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

	mux := weby.NewServeMux()
	mux.Use(webyMiddleware.RequestID)
	mux.Use(webyMiddleware.WrapResponse)
	mux.Use(webyMiddleware.Logger(logger))

	apiMux := http.NewServeMux()
	server.API(apiMux, db, contentStore, logger)
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	logger.Info("starting API server", "port", config.Port())
	serveErr := http.ListenAndServe(":"+config.Port(), mux)
	logging.Fatal(logger, "server shutting down", "err", serveErr)
}
