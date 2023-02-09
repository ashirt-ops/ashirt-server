// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"errors"
	"net/http"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/config/confighelpers"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/server"
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

	mux := http.NewServeMux()

	mux.Handle("/api/", server.API(
		db, contentStore, logger,
	))

	logger.Log("msg", "starting API server", "port", config.Port())
	serveErr := http.ListenAndServe(":"+config.Port(), mux)
	logging.Fatal(logger, "msg", "server shutting down", "err", serveErr)
}
