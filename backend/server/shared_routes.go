package server

import (
	"bytes"
	"encoding/json"
	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func bindSharedRoutes(r chi.Router, db *database.Connection, contentStore contentstore.Store) {
	route(r, "GET", "/operations", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListOperations(r.Context(), db)
	}))

	route(r, "POST", "/operations", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateOperationInput{
			Slug:    dr.FromBody("slug").Required().AsString(),
			Name:    dr.FromBody("name").Required().AsString(),
			OwnerID: middleware.UserID(r.Context()),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateOperation(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListTagsForOperationInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		return services.ListTagsForOperation(r.Context(), db, i)
	}))

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
		if s3Store, ok := contentStore.(*contentstore.S3Store); ok && evidence.ContentType == "image" {
			urlData, err := services.SendURLData(r.Context(), db, s3Store, i)
			if err != nil {
				return nil, backend.WrapError("Unable to get s3 URL", err)
			}
			jsonifiedData, err := json.Marshal(urlData)
			if err != nil {
				return nil, backend.WrapError("Unable to send s3 URL", err)
			}
			return bytes.NewReader(jsonifiedData), nil
		}
		if i.LoadPreview {
			return evidence.Preview, nil
		}
		return evidence.Media, nil
	}))

	route(r, "POST", "/operations/{operation_slug}/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateTagInput{
			Name:          dr.FromBody("name").Required().AsString(),
			ColorName:     dr.FromBody("colorName").AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateTag(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/evidence", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		timelineFilters, err := helpers.ParseTimelineQuery(dr.FromQuery("query").AsString())
		if err != nil {
			return nil, err
		}

		i := services.ListEvidenceForOperationInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Filters:       timelineFilters,
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListEvidenceForOperation(r.Context(), db, contentStore, i)
	}))

}
