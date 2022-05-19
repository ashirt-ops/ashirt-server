// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package dtos

import (
	"time"

	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/servicetypes/evidencemetadata"
)

type APIKey struct {
	AccessKey string     `json:"accessKey"`
	SecretKey []byte     `json:"secretKey"`
	LastAuth  *time.Time `json:"lastAuth"`
}

type Evidence struct {
	UUID        string             `json:"uuid"`
	Description string             `json:"description"`
	OccurredAt  time.Time          `json:"occurredAt"`
	Operator    User               `json:"operator"`
	Tags        []Tag              `json:"tags"`
	Metadata    []EvidenceMetadata `json:"metadata"` // we probably don't need this
	ContentType string             `json:"contentType"`
}

type EvidenceMetadata struct {
	Body       string                   `json:"body"`
	Source     string                   `json:"source"`
	Status     *evidencemetadata.Status `json:"status"`
	CanProcess *bool                    `json:"canProcess"`
}

type Finding struct {
	UUID          string     `json:"uuid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Operators     []User     `json:"operators"`
	ReadyToReport bool       `json:"readyToReport"`
	TicketLink    *string    `json:"ticketLink"`
	Tags          []Tag      `json:"tags"`
	NumEvidence   int        `json:"numEvidence"`
	Category      string     `json:"category"`
	OccurredFrom  *time.Time `json:"occurredFrom"`
	OccurredTo    *time.Time `json:"occurredTo"`
}

type Operation struct {
	Slug     string                 `json:"slug"`
	Name     string                 `json:"name"`
	NumUsers int                    `json:"numUsers"`
	Status   models.OperationStatus `json:"status"`
}

type Query struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
	Type  string `json:"type"`
}

type Tag struct {
	ID        int64  `json:"id"`
	ColorName string `json:"colorName"`
	Name      string `json:"name"`
}

type DefaultTag Tag

type TagWithUsage struct {
	Tag
	EvidenceCount int64 `json:"evidenceCount"`
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Slug      string `json:"slug"`
}

type UserOwnView struct {
	User
	Email          string               `json:"email"`
	Admin          bool                 `json:"admin"`
	Authentication []AuthenticationInfo `json:"authSchemes"`
	Headless       bool                 `json:"headless"`
}

type AuthenticationInfo struct {
	UserKey        string               `json:"userKey"`
	AuthSchemeName *string              `json:"schemeName,omitempty"`
	AuthSchemeCode string               `json:"schemeCode"`
	AuthSchemeType string               `json:"schemeType"`
	AuthLogin      *time.Time           `json:"lastLogin"`
	AuthDetails    *SupportedAuthScheme `json:"authDetails"`
}

type UserAdminView struct {
	User
	Email         string   `json:"email"`
	Admin         bool     `json:"admin,omitempty"`
	Headless      bool     `json:"headless"`
	Disabled      bool     `json:"disabled"`
	Deleted       bool     `json:"deleted"`
	UsesLocalTOTP bool     `json:"hasLocalTotp"`
	AuthSchemes   []string `json:"authSchemes"`
}

type UserOperationRole struct {
	User User                 `json:"user"`
	Role policy.OperationRole `json:"role"`
}

type PaginationWrapper struct {
	Content    interface{} `json:"content"`
	PageNumber int64       `json:"page"`
	PageSize   int64       `json:"pageSize"`
	TotalCount int64       `json:"totalCount"`
	TotalPages int64       `json:"totalPages"`
}

type DetailedAuthenticationInfo struct {
	AuthSchemeName  string     `json:"schemeName"`
	AuthSchemeCode  string     `json:"schemeCode"`
	AuthSchemeType  string     `json:"schemeType"`
	AuthSchemeFlags []string   `json:"schemeFlags"`
	UserCount       int64      `json:"userCount"`
	UniqueUserCount int64      `json:"uniqueUserCount"`
	LastUsed        *time.Time `json:"lastUsed"`
	Labels          []string   `json:"labels"`
}

type SupportedAuthScheme struct {
	SchemeName  string   `json:"schemeName"`
	SchemeCode  string   `json:"schemeCode"`
	SchemeFlags []string `json:"schemeFlags"`
	SchemeType  string   `json:"schemeType"`
}

type TagDifference struct {
	Included []TagPair `json:"included"`
	Excluded []Tag     `json:"excluded"`
}

type TagPair struct {
	SourceTag      Tag `json:"sourceTag"`
	DestinationTag Tag `json:"destinationTag"`
}

type TagByEvidenceDate struct {
	Tag
	UsageDates []time.Time `json:"usages"`
}

type CheckConnection struct {
	Ok bool `json:"ok"`
}

type FindingCategory struct {
	ID         int64  `json:"id"`
	Category   string `json:"category"`
	Deleted    bool   `json:"deleted"`
	UsageCount int64  `json:"usageCount"`
}

type NewUserCreatedByAdmin struct {
	TemporaryPassword string `json:"temporaryPassword"`
}

type CreateUserOutput struct {
	RealSlug string `json:"slug"`
	UserID   int64  `json:"-"` // don't transmit the userid
}

type ServiceWorker struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Config  string `json:"config"`
	Deleted bool   `json:"deleted"`
}

type ServiceWorkerTestOutput struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Live    bool   `json:"live"`
	Message string `json:"message"`
}

type ActiveServiceWorker struct {
	Name    string `json:"name"`
}
