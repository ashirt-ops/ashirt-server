// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	recoveryConsts "github.com/ashirt-ops/ashirt-server/backend/authschemes/recoveryauth/constants"
	"github.com/ashirt-ops/ashirt-server/backend/config"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/ashirt-ops/ashirt-server/backend/session"
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

func Web(r chi.Router, db *database.Connection, contentStore contentstore.Store, config *WebConfig) {
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

	r.Handle("/metrics", promhttp.Handler())
	r.Group(func(r chi.Router) {
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
			r.Route("/auth/"+scheme.Name(), func(r chi.Router) {
				scheme.BindRoutes(r.(chi.Router), authschemes.MakeAuthBridge(db, sessionStore, scheme.Name(), scheme.Type()))
			})
			supportedAuthSchemes[i] = dtos.SupportedAuthScheme{
				SchemeName:  scheme.FriendlyName(),
				SchemeCode:  scheme.Name(),
				SchemeFlags: scheme.Flags(),
				SchemeType:  scheme.Type(),
			}
		}
		authsWithOutRecovery := make([]dtos.SupportedAuthScheme, 0, len(supportedAuthSchemes)-1)

		// recovery is a special authentication that we kind of want to hide/separate from the other auth schemes
		// so, we filter it out here
		for _, auth := range supportedAuthSchemes {
			if auth.SchemeCode != recoveryConsts.Code {
				authsWithOutRecovery = append(authsWithOutRecovery, auth)
			}
		}

		bindSharedRoutes(r, db, contentStore)
		bindWebRoutes(r, db, contentStore, sessionStore, &authsWithOutRecovery)
	})
}

func bindWebRoutes(r chi.Router, db *database.Connection, contentStore contentstore.Store, sessionStore *session.Store, supportedAuthSchemes *[]dtos.SupportedAuthScheme) {
	route(r, "POST", "/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonHandler(func(r *http.Request) (interface{}, error) {
			err := sessionStore.Delete(w, r)
			if err != nil {
				return nil, backend.WrapError("Unable to delete session", err)
			}
			return nil, nil
		}).ServeHTTP(w, r)
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

	route(r, "GET", "/usergroups", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListUserGroupsInput{
			Query:          dr.FromQuery("query").Required().AsString(),
			IncludeDeleted: dr.FromQuery("includeDeleted").OrDefault(false).AsBool(),
			OperationSlug:  dr.FromQuery("operationSlug").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListUserGroups(r.Context(), db, i)
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

	route(r, "GET", "/admin/usergroups", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListUserGroupsForAdminInput{
			UserGroupFilter: services.ParseRequestQueryUserGroupFilter(dr),
			IncludeDeleted:  dr.FromQuery("deleted").OrDefault(false).AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListUserGroupsForAdmin(r.Context(), db, i)
	}))

	route(r, "POST", "/admin/usergroups", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateUserGroupInput{
			Slug:      dr.FromBody("slug").Required().AsString(),
			Name:      dr.FromBody("name").Required().AsString(),
			UserSlugs: dr.FromBody("userSlugs").Required().AsStringSlice(),
		}

		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateUserGroup(r.Context(), db, i)
	}))

	route(r, "PUT", "/admin/usergroups/{group_slug}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ModifyUserGroupInput{
			Name:          dr.FromBody("newName").AsString(),
			UsersToAdd:    dr.FromBody("userSlugsToAdd").AsStringSlice(),
			UsersToRemove: dr.FromBody("userSlugsToRemove").AsStringSlice(),
			Slug:          dr.FromURL("group_slug").Required().AsString(),
		}

		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ModifyUserGroup(r.Context(), db, i)
	}))

	route(r, "DELETE", "/admin/usergroups/{group_slug}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		groupSlug := dr.FromURL("group_slug").AsString()

		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteUserGroup(r.Context(), db, groupSlug)
	}))

	route(r, "GET", "/auths", jsonHandler(func(r *http.Request) (interface{}, error) {
		return supportedAuthSchemes, nil
	}))

	route(r, "GET", "/flags", jsonHandler(func(r *http.Request) (interface{}, error) {
		return dtos.Flags{Flags: config.Flags()}, nil
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

	route(r, "GET", "/admin/operations", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListOperationsForAdmin(r.Context(), db)
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

	route(r, "GET", "/operations/{operation_slug}/usergroups", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ListUserGroupsForOperationInput{
			OperationSlug:   dr.FromURL("operation_slug").Required().AsString(),
			UserGroupFilter: services.ParseRequestQueryUserGroupFilter(dr),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}

		return services.ListUserGroupsForOperation(r.Context(), db, i)
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

	route(r, "PATCH", "/operations/{operation_slug}/usergroups", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.SetUserGroupOperationRoleInput{
			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
			UserGroupSlug: dr.FromBody("userGroupSlug").Required().AsString(),
			Role:          policy.OperationRole(dr.FromBody("role").Required().AsString()),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.SetUserGroupOperationRole(r.Context(), db, i)
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
		return services.ListEvidenceForOperation(r.Context(), db, contentStore, i)
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

	var handler http.Handler

	// if s3Store, ok := contentStore.(*contentstore.S3Store); ok {
	// 	handler = jsonHandler(func(r *http.Request) (interface{}, error) {
	// 		dr := dissectNoBodyRequest(r)
	// 		i := services.ReadEvidenceInput{
	// 			EvidenceUUID:  dr.FromURL("evidence_uuid").Required().AsString(),
	// 			OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
	// 			LoadPreview:   dr.FromURL("type").AsString() == "preview",
	// 			LoadMedia:     dr.FromURL("type").AsString() == "media",
	// 		}

	// 		url, err := services.SendURL(r.Context(), db, s3Store, i)
	// 		if err != nil {
	// 			return nil, backend.WrapError("Unable to obtain image URL", err)
	// 		}
	// 		return url, nil
	// 	})
	// } else {
	handler = mediaHandler(func(r *http.Request) (io.Reader, error) {
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
		// TODO TN - maybe we want all stuff sent out of band?
		if s3Store, ok := contentStore.(*contentstore.S3Store); ok && evidence.ContentType == "image" {
			url, _ := services.SendURL2(r.Context(), db, s3Store, i)
			fmt.Println("___*url", *url)

			reader := strings.NewReader(*url)
			// var genericReader io.Reader
			// genericReader = reader

			// Use the io.Reader as needed
			// buffer := make([]byte, 32)
			// // result := ""
			// for {
			// 	_, err := genericReader.Read(buffer)
			// 	// n, err := reader.Read(buffer)
			// 	if err == io.EOF {
			// 		break
			// 	}
			// 	if err != nil {
			// 		fmt.Println("Error:", err)
			// 		break
			// 	}
			// 	// result += string(buffer[:n])
			// }

			// Print the result to verify
			// fmt.Println("Content from io.Reader:", result, *url)
			// // Compare the result with the original string
			// if result == *url {
			// 	fmt.Println("Content matches the original string.")
			// } else {
			// 	fmt.Println("Content does not match the original string.")
			// }

			return reader, nil
			// return genericReader, nil
		}

		if i.LoadPreview {
			return evidence.Preview, nil
		}
		return evidence.Media, nil
	})
	// }

	route(r, "GET", "/operations/{operation_slug}/evidence/{evidence_uuid}/{type:media|preview}", handler)

	route(r, "GET", "/operations/{operation_slug}/evidence/{evidence_uuid}/metadata", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.ReadEvidenceMetadataInput{
			OperationSlug: dr.FromURL("operation_slug").AsString(),
			EvidenceUUID:  dr.FromURL("evidence_uuid").AsString(),
		}
		return services.ReadEvidenceMetadata(r.Context(), db, i)
	}))

	route(r, "POST", "/operations/{operation_slug}/evidence/{evidence_uuid}/metadata", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.EditEvidenceMetadataInput{
			OperationSlug: dr.FromURL("operation_slug").AsString(),
			EvidenceUUID:  dr.FromURL("evidence_uuid").AsString(),
			Source:        dr.FromBody("source").Required().AsString(),
			Body:          dr.FromBody("body").Required().AsString(),
		}
		return nil, services.CreateEvidenceMetadata(r.Context(), db, i)
	}))

	route(r, "PUT", "/operations/{operation_slug}/evidence/{evidence_uuid}/metadata", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.EditEvidenceMetadataInput{
			OperationSlug: dr.FromURL("operation_slug").AsString(),
			EvidenceUUID:  dr.FromURL("evidence_uuid").AsString(),
			Source:        dr.FromBody("source").Required().AsString(),
			Body:          dr.FromBody("body").Required().AsString(),
		}
		return nil, services.UpdateEvidenceMetadata(r.Context(), db, i)
	}))

	route(r, "PUT", "/operations/{operation_slug}/evidence/{evidence_uuid}/metadata/{service_name}/run", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.RunServiceWorkerInput{
			OperationSlug: dr.FromURL("operation_slug").AsString(),
			EvidenceUUID:  dr.FromURL("evidence_uuid").AsString(),
			WorkerName:    dr.FromURL("service_name").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.RunServiceWorker(r.Context(), db, i)
	}))

	route(r, "PUT", "/operations/{operation_slug}/evidence/{evidence_uuid}/metadata/run", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.RunServiceWorkerInput{
			OperationSlug: dr.FromURL("operation_slug").AsString(),
			EvidenceUUID:  dr.FromURL("evidence_uuid").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.RunServiceWorker(r.Context(), db, i)
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

	route(r, "PUT", "/operations/{operation_slug}/queries", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpsertQueryInput{
			CreateQueryInput: services.CreateQueryInput{
				OperationSlug: dr.FromURL("operation_slug").Required().AsString(),
				Name:          dr.FromBody("name").Required().AsString(),
				Query:         dr.FromBody("query").Required().AsString(),
				Type:          dr.FromBody("type").Required().AsString(),
			},
			ReplaceName: dr.FromBody("replaceName").OrDefault(false).AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.UpsertQuery(r.Context(), db, i)
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

	route(r, "GET", "/admin/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListDefaultTags(r.Context(), db)
	}))

	route(r, "POST", "/admin/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateDefaultTagInput{
			Name:      dr.FromBody("name").Required().AsString(),
			ColorName: dr.FromBody("colorName").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateDefaultTag(r.Context(), db, i)
	}))

	route(r, "POST", "/admin/merge/tags", jsonHandler(func(r *http.Request) (interface{}, error) {
		var tags []services.CreateDefaultTagInput
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &tags); err != nil {
			return nil, err
		}
		return nil, services.MergeDefaultTags(r.Context(), db, tags)
	}))

	route(r, "PUT", "/admin/tags/{tag_id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateDefaultTagInput{
			ID:        dr.FromURL("tag_id").Required().AsInt64(),
			Name:      dr.FromBody("name").Required().AsString(),
			ColorName: dr.FromBody("colorName").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateDefaultTag(r.Context(), db, i)
	}))

	route(r, "DELETE", "/admin/tags/{tag_id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteDefaultTagInput{
			ID: dr.FromURL("tag_id").Required().AsInt64(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteDefaultTag(r.Context(), db, i)
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

	route(r, "GET", "/findings/categories", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		includeDeleted := dr.FromQuery("includeDeleted").OrDefault(false).AsBool()

		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.ListFindingCategories(r.Context(), db, includeDeleted)
	}))

	route(r, "POST", "/findings/category", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		category := dr.FromBody("category").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateFindingCategory(r.Context(), db, category)
	}))

	route(r, "DELETE", "/findings/category/{id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteFindingCategoryInput{
			FindingCategoryID: dr.FromURL("id").AsInt64(),
			DoDelete:          dr.FromBody("delete").Required().AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteFindingCategory(r.Context(), db, i)
	}))

	route(r, "PUT", "/findings/category/{id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateFindingCategoryInput{
			Category: dr.FromBody("category").Required().AsString(),
			ID:       dr.FromURL("id").AsInt64(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateFindingCategory(r.Context(), db, i)
	}))

	bindServiceWorkerRoutes(r, db)
}

func bindServiceWorkerRoutes(r chi.Router, db *database.Connection) {
	route(r, "GET", "/admin/services", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListServiceWorker(r.Context(), db)
	}))

	route(r, "GET", "/services", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListActiveServices(r.Context(), db)
	}))

	route(r, "POST", "/admin/services", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateServiceWorkerInput{
			Name:   dr.FromBody("name").Required().AsString(),
			Config: dr.FromBody("config").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.CreateServiceWorker(r.Context(), db, i)
	}))

	route(r, "PUT", "/admin/services/{id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateServiceWorkerInput{
			ID:     dr.FromURL("id").AsInt64(),
			Name:   dr.FromBody("name").Required().AsString(),
			Config: dr.FromBody("config").Required().AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateServiceWorker(r.Context(), db, i)
	}))

	route(r, "DELETE", "/admin/services/{id}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.DeleteServiceWorkerInput{
			ID:       dr.FromURL("id").AsInt64(),
			DoDelete: dr.FromBody("delete").Required().AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteServiceWorker(r.Context(), db, i)
	}))

	route(r, "GET", "/admin/services/{id}/test", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		workerID := dr.FromURL("id").AsInt64()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.TestServiceWorker(r.Context(), db, workerID)
	}))

	route(r, "PUT", "/operations/{operation_slug}/metadata/run", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.BatchRunServiceWorkerInput{
			OperationSlug: dr.FromURL("operation_slug").AsString(),
			EvidenceUUIDs: dr.FromBody("evidenceUuids").Required().AsStringSlice(),
			WorkerNames:   dr.FromBody("workers").Required().AsStringSlice(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.BatchRunServiceWorker(r.Context(), db, i)
	}))

	route(r, "POST", "/operations/{operation_slug}/favorite", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.SetFavoriteInput{
			OperationSlug: dr.FromURL("operation_slug").AsString(),
			IsFavorite:    dr.FromBody("favorite").Required().AsBool(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.SetFavoriteOperation(r.Context(), db, i)
	}))

	route(r, "GET", "/global-vars", jsonHandler(func(r *http.Request) (interface{}, error) {
		return services.ListGlobalVars(r.Context(), db)
	}))

	route(r, "POST", "/global-vars", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.CreateGlobalVarInput{
			Name:  dr.FromBody("name").Required().AsString(),
			Value: dr.FromBody("value").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return services.CreateGlobalVar(r.Context(), db, i)
	}))

	route(r, "PUT", "/global-vars/{name}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		i := services.UpdateGlobalVarInput{
			Name:    dr.FromURL("name").Required().AsString(),
			Value:   dr.FromBody("value").AsString(),
			NewName: dr.FromBody("newName").AsString(),
		}
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.UpdateGlobalVar(r.Context(), db, i)
	}))

	route(r, "DELETE", "/global-vars/{name}", jsonHandler(func(r *http.Request) (interface{}, error) {
		dr := dissectJSONRequest(r)
		name := dr.FromURL("name").Required().AsString()
		if dr.Error != nil {
			return nil, dr.Error
		}
		return nil, services.DeleteGlobalVar(r.Context(), db, name)
	}))
}
