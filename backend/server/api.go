// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/logging"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/server/middleware"
	"github.com/theparanoids/ashirt/backend/services"

	sq "github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
)

var operationIDRegex = regexp.MustCompile(`\d+`)

func API(db *database.Connection, contentStore contentstore.Store, logger logging.Logger) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware.LogRequests(logger))
	r.Use(middleware.AuthenticateAppAndInjectCtx(db))

	bindAPIRoutes(r, db, contentStore)
	return r
}

// Temporary for now since the api uses ids still but the services have been swapped over to using slugs
// Once the screenshot client is updated to use slugs this function can go away:
func operationIDToSlug(db *database.Connection, operationID int64) string {
	var operation models.Operation
	err := db.Get(&operation, sq.Select("slug").From("operations").Where(sq.Eq{"id": operationID}))
	if err != nil {
		// If we can't find it the service will return a NotFoundErr for
		// an empty string operationslug so no need to handle errors here
		return ""
	}
	return operation.Slug
}

func bindAPIRoutes(r *mux.Router, db *database.Connection, contentStore contentstore.Store) {
	route(r, "GET", "/api/operations", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListOperations(r.Context(), db)
	}))

	route(r, "GET", "/api/operations/{operation_id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		operationID := dr.FromURL("operation_id").Required().AsInt64()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ReadOperation(r.Context(), db, operationIDToSlug(db, operationID))
	}))

	route(r, "POST", "/api/operations/{operation_id}/evidence", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectFormRequest(r)
		opSlug := dr.FromURL("operation_id").Required().AsString()
		opSlug = maybeOperationIDToSlug(db, opSlug)
		i := services.CreateEvidenceInput{
			Description:   dr.FromBody("notes").Required().AsString(),
			Content:       dr.FromFile("file"),
			ContentType:   dr.FromBody("contentType").OrDefault("image").AsString(),
			OccurredAt:    dr.FromBody("occurred_at").OrDefault(time.Now()).AsUnixTime(),
			OperationSlug: opSlug,
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

	route(r, "GET", "/api/operations/{operation_slug}/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		opSlug := dr.FromURL("operation_slug").Required().AsString()
		opSlug = maybeOperationIDToSlug(db, opSlug)

		i := services.ListTagsForOperationInput{
			OperationSlug: opSlug,
		}
		return services.ListTagsForOperation(r.Context(), db, i)
	}))

	route(r, "POST", "/api/operations/{operation_slug}/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		opSlug := dr.FromURL("operation_slug").Required().AsString()

		opSlug = maybeOperationIDToSlug(db, opSlug)

		i := services.CreateTagInput{
			Name:          dr.FromBody("name").Required().AsString(),
			ColorName:     dr.FromBody("colorName").AsString(),
			OperationSlug: opSlug,
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateTag(r.Context(), db, i)
	}))
}

func maybeOperationIDToSlug(db *database.Connection, slugOrID string) string {
	if operationIDRegex.MatchString(slugOrID) {
		// id to slug
		if val, err := strconv.ParseInt(slugOrID, 10, 64); err == nil {
			return operationIDToSlug(db, val)
		}
	}
	return slugOrID // must actually be a slug
}
