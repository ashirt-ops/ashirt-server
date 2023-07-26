// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
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
		bindSharedRoutes(r, db, contentStore)
		bindAPIRoutes(r, db, contentStore)
	})
}

func bindAPIRoutes(r chi.Router, db *database.Connection, contentStore contentstore.Store) {
	route(r, "GET", "/operations/{operation_slug}/evidence/{evidence_uuid}/{type:media|preview}", mediaHandler(func(r *http.Request) (io.Reader, error) {
		dr := dissectNoBodyRequest(r)
		i := services.ReadEvidenceInput{
			EvidenceUUID:  dr.FromURL("evidence_uuid").Required().AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			LoadPreview:   dr.FromURL("type").AsString() == "preview",
			LoadMedia:     dr.FromURL("type").AsString() == "media",
		}

		evidence, err := services.ReadEvidence(r.Context(), db, contentStore, i)
		if err != nil {
			return nil, backend.WrapError("Unable to read evidence", err)
		}

		if i.LoadPreview {
			return evidence.Preview, nil
		}
		return evidence.Media, nil
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

	route(r, "PUT", "/operations/{operation_slug}/evidence/{evidence_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectFormRequest(r)
		i := services.UpdateEvidenceInput{
			EvidenceUUID:  dr.FromURL("evidence_uuid").Required().AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Description:   dr.FromBody("notes").AsStringPtr(),
			Content:       dr.FromFile("file"),
		}
		tagsToAddJSON := dr.FromBody("tagsToAdd").OrDefault("[]").AsString()
		tagsToRemoveJSON := dr.FromBody("tagsToRemove").OrDefault("[]").AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		if err := json.Unmarshal([]byte(tagsToAddJSON), &i.TagsToAdd); err != nil {
			return nil, backend.BadInputErr(err, "tagsToAdd must be a json array of ints")
		}
		if err := json.Unmarshal([]byte(tagsToRemoveJSON), &i.TagsToRemove); err != nil {
			return nil, backend.BadInputErr(err, "tagsToRemove must be a json array of ints")
		}
		return nil, services.UpdateEvidence(r.Context(), db, contentStore, i)
	}))

	route(r, "PUT", "/operations/{operation_slug}/evidence/{evidence_uuid}/metadata", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpsertEvidenceMetadataInput{
			EditEvidenceMetadataInput: services.EditEvidenceMetadataInput{
				OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
				EvidenceUUID:  dr.FromURL("evidence_uuid").Required().AsString(),
				Source:        dr.FromBody("source").Required().AsString(),
				Body:          dr.FromBody("body").Required().AsString(),
			},
			Status:     dr.FromBody("status").Required().AsString(),
			Message:    dr.FromBody("message").AsStringPtr(),
			CanProcess: dr.FromBody("canProcess").AsBoolPtr(),
		}
		return nil, services.UpsertEvidenceMetadata(r.Context(), db, i)
	}))
}
