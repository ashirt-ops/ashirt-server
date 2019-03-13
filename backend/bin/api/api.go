// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"net/http"

	"github.com/theparanoids/ashirt/backend/config"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/logging"
	"github.com/theparanoids/ashirt/backend/server"
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

	contentStore, err := contentstore.NewS3Store(config.ImageStoreBucketName(), config.ImageStoreRegion())
	if err != nil {
		logging.Fatal(logger, "msg", "store setup error", "error", err)
	}

	http.Handle("/api/", server.API(
		db, contentStore, logger,
	))

	logger.Log("msg", "starting API server", "port", config.Port())
	serveErr := http.ListenAndServe(":"+config.Port(), nil)
	logging.Fatal(logger, "msg", "server shutting down", "err", serveErr)
}
