// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/theparanoids/ashirt-server/backend/dtos"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	recoveryConsts "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"github.com/theparanoids/ashirt-server/backend/services"
	"github.com/theparanoids/ashirt-server/backend/session"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type WebConfig struct {
	DBConnection     *database.Connection
	AuthSchemes      []authschemes.AuthScheme
	CSRFAuthKey      []byte
	SessionStoreKey  []byte
	UseSecureCookies bool
	Logger           logging.Logger
}

func (c *WebConfig) validate() error {
	if c.Logger == nil {
		fmt.Println(`error="Logger not set" action="Using NopLogger"`)
		c.Logger = logging.NewNopLogger()
	}
	if len(c.CSRFAuthKey) < 32 {
		return errors.New("CSRFAuthKey must be 32 bytes or longer")
	}
	if len(c.SessionStoreKey) < 32 {
		return errors.New("SessionStoreKey must be 32 bytes or longer")
	}
	if !c.UseSecureCookies {
		c.Logger.Log("msg", "Config Warning: cookies not using secure flag")
	}
	return nil
}

func Web(db *database.Connection, contentStore contentstore.Store, config *WebConfig) http.Handler {
	if err := config.validate(); err != nil {
		panic(err)
	}
	sessionStore, err := session.NewStore(db, session.StoreOptions{
		SessionDuration:  30 * 24 * time.Hour,
		UseSecureCookies: config.UseSecureCookies,
		Key:              config.SessionStoreKey,
	})
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Use(middleware.LogRequests(config.Logger))
	r.Use(csrf.Protect(config.CSRFAuthKey,
		csrf.Secure(config.UseSecureCookies),
		csrf.Path("/"),
		csrf.ErrorHandler(jsonHandler(func(r *http.Request) (interface{}, error) {
			return nil, backend.CSRFErr(csrf.FailureReason(r))
		}))))
	r.Use(middleware.InjectCSRFTokenHeader())
	r.Use(middleware.AuthenticateUserAndInjectCtx(db, sessionStore))

	supportedAuthSchemes := make([]dtos.SupportedAuthScheme, len(config.AuthSchemes))
	for i, scheme := range config.AuthSchemes {
		authRouter := r.PathPrefix("/auth/" + scheme.Name()).Subrouter()
		scheme.BindRoutes(authRouter, authschemes.MakeAuthBridge(db, sessionStore, scheme.Name()))
		supportedAuthSchemes[i] = dtos.SupportedAuthScheme{SchemeName: scheme.FriendlyName(), SchemeCode: scheme.Name()}
	}
	authsWithOutRecovery := make([]dtos.SupportedAuthScheme, 0, len(supportedAuthSchemes)-1)

	// recovery is a special authentication that we kind of want to hide/separate from the other auth schemes
	// so, we filter it out here
	for _, auth := range supportedAuthSchemes {
		if auth.SchemeCode != recoveryConsts.Code {
			authsWithOutRecovery = append(authsWithOutRecovery, auth)
		}
	}

	bindWebRoutes(r, db, contentStore, sessionStore, &authsWithOutRecovery)
	return r
}

func bindWebRoutes(r *mux.Router, db *database.Connection, contentStore contentstore.Store, sessionStore *session.Store, supportedAuthSchemes *[]dtos.SupportedAuthScheme) {
	route(r, "POST", "/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonHandler(func(r *http.Request) (interface{}, error) {
			err := sessionStore.Delete(w, r)
			if err != nil {
				return nil, backend.WrapError("Unable to delete session", err)
			}
			return nil, nil
		}).ServeHTTP(w, r)
	}))

	route(r, "GET", "/health", jsonHandler(func(r *http.Request) (interface{}, error) {
		return nil, nil
	}))

	route(r, "GET", "/user", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		slug := dr.FromQuery("userSlug").AsString()

		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ReadUser(r.Context(), db, slug, supportedAuthSchemes)
	}))

	route(r, "GET", "/users", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListUsersInput{
			Query:          dr.FromQuery("query").Required().AsString(),
			IncludeDeleted: dr.FromQuery("includeDeleted").OrDefault(false).AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListUsers(r.Context(), db, i)
	}))

	route(r, "GET", "/admin/users", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListUsersForAdminInput{
			UserFilter:     services.ParseRequestQueryUserFilter(dr),
			Pagination:     services.ParseRequestQueryPagination(dr, 10),
			IncludeDeleted: dr.FromQuery("deleted").OrDefault(false).AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListUsersForAdmin(r.Context(), db, i)
	}))

	route(r, "DELETE", "/admin/user/{userSlug}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := dr.FromURL("userSlug").AsString()

		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteUser(r.Context(), db, i)
	}))

	route(r, "POST", "/admin/user/headless", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)

		i := services.CreateUserInput{
			FirstName: dr.FromBody("firstName").AsString(),
			LastName:  dr.FromBody("lastName").AsString(),
			Email:     dr.FromBody("email").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		i.Slug = i.FirstName + "." + i.LastName

		return services.CreateHeadlessUser(r.Context(), db, i)
	}))

	route(r, "POST", "/admin/{userSlug}/flags", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.SetUserFlagsInput{
			Slug:     dr.FromURL("userSlug").AsString(),
			Disabled: dr.FromBody("disabled").AsBoolPtr(),
			Admin:    dr.FromBody("admin").AsBoolPtr(),
		}

		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.SetUserFlags(r.Context(), db, i)
	}))

	route(r, "GET", "/auths", jsonHandler(func(r *http.Request) (interface{}, error) {
		return supportedAuthSchemes, nil
	}))

	route(r, "GET", "/auths/breakdown", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListAuthDetails(r.Context(), db, supportedAuthSchemes)
	}))

	route(r, "DELETE", "/auths/{schemeCode}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		schemeCode := dr.FromURL("schemeCode").AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteAuthSchemeUsers(r.Context(), db, schemeCode)
	}))

	route(r, "GET", "/operations", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListOperations(r.Context(), db)
	}))

	route(r, "GET", "/admin/operations", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListOperationsForAdmin(r.Context(), db)
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

	route(r, "DELETE", "/operations/{operation_slug}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		operationSlug := dr.FromURL("operation_slug").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}

		return nil, services.DeleteOperation(r.Context(), db, contentStore, operationSlug)
	}))

	route(r, "GET", "/operations/{operation_slug}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		operationSlug := dr.FromURL("operation_slug").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}

		return services.ReadOperation(r.Context(), db, operationSlug)
	}))

	route(r, "PUT", "/operations/{operation_slug}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateOperationInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Name:          dr.FromBody("name").Required().AsString(),
			Status:        models.OperationStatus(dr.FromBody("status").OrDefault(int64(models.OperationStatusPlanning)).AsInt64()),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateOperation(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/users", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListUsersForOperationInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			UserFilter:    services.ParseRequestQueryUserFilter(dr),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}

		return services.ListUsersForOperation(r.Context(), db, i)
	}))

	route(r, "PATCH", "/operations/{operation_slug}/users", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.SetUserOperationRoleInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			UserSlug:      dr.FromBody("userSlug").Required().AsString(),
			Role:          policy.OperationRole(dr.FromBody("role").Required().AsString()),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.SetUserOperationRole(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/findings", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		timelineFilters, err := helpers.ParseTimelineQuery(dr.FromQuery("query").AsString())
		if err != nil {
			return nil, backend.WrapError("Unable to parse findings query", err)
		}

		i := services.ListFindingsForOperationInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Filters:       timelineFilters,
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListFindingsForOperation(r.Context(), db, i)
	}))

	route(r, "POST", "/operations/{operation_slug}/findings", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateFindingInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Category:      dr.FromBody("category").Required().AsString(),
			Title:         dr.FromBody("title").Required().AsString(),
			Description:   dr.FromBody("description").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateFinding(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/findings/{finding_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ReadFindingInput{
			FindingUUID:   dr.FromURL("finding_uuid").Required().AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ReadFinding(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/findings/{finding_uuid}/evidence", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListEvidenceForFindingInput{
			FindingUUID:   dr.FromURL("finding_uuid").Required().AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListEvidenceForFinding(r.Context(), db, i)
	}))

	route(r, "PUT", "/operations/{operation_slug}/findings/{finding_uuid}/evidence", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.AddEvidenceToFindingInput{
			OperationSlug:    dr.FromURL("operation_slug").Required().AsString(),
			FindingUUID:      dr.FromURL("finding_uuid").Required().AsString(),
			EvidenceToAdd:    dr.FromBody("evidenceToAdd").Required().AsStringSlice(),
			EvidenceToRemove: dr.FromBody("evidenceToRemove").Required().AsStringSlice(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.AddEvidenceToFinding(r.Context(), db, i)
	}))

	route(r, "PUT", "/operations/{operation_slug}/findings/{finding_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateFindingInput{
			FindingUUID:   dr.FromURL("finding_uuid").Required().AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Category:      dr.FromBody("category").Required().AsString(),
			Title:         dr.FromBody("title").AsString(),
			Description:   dr.FromBody("description").AsString(),
			TicketLink:    dr.FromBody("ticketLink").AsStringPtr(),
			ReadyToReport: dr.FromBody("readyToReport").Required().AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateFinding(r.Context(), db, i)
	}))

	route(r, "DELETE", "/operations/{operation_slug}/findings/{finding_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteFindingInput{
			FindingUUID:   dr.FromURL("finding_uuid").Required().AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteFinding(r.Context(), db, i)
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
		return services.ListEvidenceForOperation(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/evidence/creators", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)

		i := services.ListEvidenceCreatorsForOperationInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListEvidenceCreatorsForOperation(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/evidence/{evidence_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ReadEvidenceInput{
			EvidenceUUID:  dr.FromURL("evidence_uuid").Required().AsString(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ReadEvidence(r.Context(), db, contentStore, i)
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

		if i.LoadPreview {
			return evidence.Preview, nil
		}
		return evidence.Media, nil
	}))

	route(r, "POST", "/operations/{operation_slug}/evidence", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectFormRequest(r)
		i := services.CreateEvidenceInput{
			Description:   dr.FromBody("description").Required().AsString(),
			Content:       dr.FromFile("content"),
			ContentType:   dr.FromBody("contentType").OrDefault("image").AsString(),
			OccurredAt:    dr.FromBody("occurredAt").OrDefault(time.Now()).AsTime(),
			OperationSlug: dr.FromURL("operation_slug").AsString(),
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
			Description:   dr.FromBody("description").AsStringPtr(),
			Content:       dr.FromFile("content"),
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

	route(r, "PUT", "/move/operations/{operation_slug}/evidence/{evidence_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.MoveEvidenceInput{
			EvidenceUUID:        dr.FromURL("evidence_uuid").Required().AsString(),
			TargetOperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			SourceOperationSlug: dr.FromBody("sourceOperationSlug").Required().AsString(),
		}

		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.MoveEvidence(r.Context(), db, i)
	}))

	route(r, "GET", "/move/operations/{operation_slug}/evidence/{evidence_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListTagDifferenceForEvidenceInput{
			ListTagsDifferenceInput: services.ListTagsDifferenceInput{
				SourceOperationSlug:      dr.FromQuery("sourceOperationSlug").Required().AsString(),
				DestinationOperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			},
			SourceEvidenceUUID: dr.FromURL("evidence_uuid").Required().AsString(),
		}
		return services.ListTagDifferenceForEvidence(r.Context(), db, i)
	}))

	route(r, "DELETE", "/operations/{operation_slug}/evidence/{evidence_uuid}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteEvidenceInput{
			EvidenceUUID:             dr.FromURL("evidence_uuid").Required().AsString(),
			OperationSlug:            dr.FromURL("operation_slug").Required().AsString(),
			DeleteAssociatedFindings: dr.FromBody("deleteAssociatedFindings").OrDefault(false).AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteEvidence(r.Context(), db, contentStore, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/queries", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		operationID := dr.FromURL("operation_slug").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListQueriesForOperation(r.Context(), db, operationID)
	}))

	route(r, "POST", "/operations/{operation_slug}/queries", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateQueryInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Name:          dr.FromBody("name").Required().AsString(),
			Query:         dr.FromBody("query").Required().AsString(),
			Type:          dr.FromBody("type").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateQuery(r.Context(), db, i)
	}))

	route(r, "PUT", "/operations/{operation_slug}/queries/{query_id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateQueryInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			ID:            dr.FromURL("query_id").Required().AsInt64(),
			Name:          dr.FromBody("name").Required().AsString(),
			Query:         dr.FromBody("query").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateQuery(r.Context(), db, i)
	}))

	route(r, "DELETE", "/operations/{operation_slug}/queries/{query_id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteQueryInput{
			ID:            dr.FromURL("query_id").Required().AsInt64(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteQuery(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListTagsForOperationInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		return services.ListTagsForOperation(r.Context(), db, i)
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

	route(r, "PUT", "/operations/{operation_slug}/tags/{tag_id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateTagInput{
			ID:            dr.FromURL("tag_id").Required().AsInt64(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			Name:          dr.FromBody("name").Required().AsString(),
			ColorName:     dr.FromBody("colorName").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateTag(r.Context(), db, i)
	}))

	route(r, "DELETE", "/operations/{operation_slug}/tags/{tag_id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteTagInput{
			ID:            dr.FromURL("tag_id").Required().AsInt64(),
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteTag(r.Context(), db, i)
	}))

	route(r, "GET", "/operations/{operation_slug}/tagsByEvidenceUsage", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListTagsByEvidenceDateInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
		}
		return services.ListTagsByEvidenceDate(r.Context(), db, i)
	}))

	route(r, "GET", "/user/apikeys", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		userSlug := dr.FromQuery("userSlug").AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListAPIKeys(r.Context(), db, userSlug)
	}))

	route(r, "POST", "/user/{userSlug}/apikeys", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		userSlug := dr.FromURL("userSlug").AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateAPIKey(r.Context(), db, userSlug)
	}))

	route(r, "DELETE", "/user/{userSlug}/apikeys/{access_key}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteAPIKeyInput{
			UserSlug:  dr.FromURL("userSlug").Required().AsString(),
			AccessKey: dr.FromURL("access_key").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteAPIKey(r.Context(), db, i)
	}))

	route(r, "POST", "/user/profile/{userSlug}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateUserProfileInput{
			UserSlug:  dr.FromURL("userSlug").AsString(),
			FirstName: dr.FromBody("firstName").Required().AsString(),
			LastName:  dr.FromBody("lastName").Required().AsString(),
			Email:     dr.FromBody("email").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateUserProfile(r.Context(), db, i)
	}))

	route(r, "DELETE", "/user/{userSlug}/scheme/{authSchemeName}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteAuthSchemeInput{
			UserSlug:   dr.FromURL("userSlug").AsString(),
			SchemeName: dr.FromURL("authSchemeName").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteAuthScheme(r.Context(), db, i)
	}))
}
