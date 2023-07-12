// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package server

import (
	"net/http"

	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func API(r chi.Router, db *database.Connection, contentStore contentstore.Store, logger logging.Logger) {
	r.Handle("/metrics", promhttp.Handler())
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthenticateAppAndInjectCtx(db))
		r.Use(middleware.LogRequests(logger))

		route(r, "GET", "/checkconnection", jsonHandler(func(r *http.Request) (interface{}, error) {
			return dtos.CheckConnection{Ok: true}, nil
		}))
		bindSharedRoutes(r, db, contentStore)
	})
}
