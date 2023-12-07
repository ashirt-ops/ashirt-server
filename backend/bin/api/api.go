// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

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
	"github.com/go-chi/chi/v5"
)

func main() {
	err := config.LoadAPIConfig()
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

	s := chi.NewRouter()

	s.Route("/api", func(r chi.Router) {
		server.API(r,
			db, contentStore, logger,
		)
	})

	logger.Log("msg", "starting API server", "port", config.Port())
	serveErr := http.ListenAndServe(":"+config.Port(), s)
	logging.Fatal(logger, "msg", "server shutting down", "err", serveErr)
}
