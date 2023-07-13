// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/ashirt-ops/ashirt-server/backend/services"
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
		route(r, "POST", "/operations/{operation_slug}/evidence", jsonHandler(func(r *http.Request) (interface{}, error) {
			dr := dissectFormRequest(r)
			i := services.CreateEvidenceInput{
				Description:   dr.FromBody("notes").Required().AsString(),
				Content:       dr.FromFile("file"),
				ContentType:   dr.FromBody("contentType").OrDefault("image").AsString(),
				OccurredAt:    dr.FromBody("occurred_at").OrDefault(time.Now()).AsUnixTime(),
				OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			}
			tagIDsJSON := dr.FromBody("tagIds").OrDefault("[]").AsString()
			if dr.Error != nil {
				return nil, dr.Error
			}
			if err := json.Unmarshal([]byte(tagIDsJSON), &i.TagIDs); err != nil {
				return nil, backend.BadInputErr(err, "tagIds must be a json array of ints")
			}
			return services.CreateEvidence(r.Context(), db, contentStore, i)
		}))
		bindSharedRoutes(r, db, contentStore)
	})
}
