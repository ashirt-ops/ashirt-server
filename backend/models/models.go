// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package models

import (
	"time"

	"github.com/theparanoids/ashirt/backend/policy"
)

// APIKey reflects the structure of the database table 'api_keys'
type APIKey struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	AccessKey string     `db:"access_key"`
	SecretKey []byte     `db:"secret_key"`
	LastAuth  *time.Time `db:"last_auth"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// Finding reflects the structure of the database table 'findings'
type Finding struct {
	ID            int64      `db:"id" json:"id"`
	UUID          string     `db:"uuid" json:"uuid"`
	OperationID   int64      `db:"operation_id" json:"operationId"`
	ReadyToReport bool       `db:"ready_to_report" json:"readyToReport"`
	TicketLink    *string    `db:"ticket_link" json:"ticketLink"`
	Category      string     `db:"category" json:"category"`
	Title         string     `db:"title" json:"title"`
	Description   string     `db:"description" json:"description"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updatedAt"`
}

// Evidence reflects the structure of the database table 'evidence'
type Evidence struct {
	ID            int64      `db:"id"              json:"id"`
	UUID          string     `db:"uuid"            json:"uuid"`
	OperationID   int64      `db:"operation_id"    json:"operationId"`
	OperatorID    int64      `db:"operator_id"     json:"operatorId"`
	Description   string     `db:"description"     json:"description"`
	ContentType   string     `db:"content_type"    json:"contentType"`
	FullImageKey  string     `db:"full_image_key"  json:"fullImageKey"`
	ThumbImageKey string     `db:"thumb_image_key" json:"thumbImageKey"`
	OccurredAt    time.Time  `db:"occurred_at"     json:"occurredAt"`
	CreatedAt     time.Time  `db:"created_at"      json:"createdAt"`
	UpdatedAt     *time.Time `db:"updated_at"      json:"updatedAt"`
}

// EvidenceFindingMap reflects the structure of the database table 'evidence_finding_map'
type EvidenceFindingMap struct {
	EvidenceID int64      `db:"evidence_id" json:"evidenceId"`
	FindingID  int64      `db:"finding_id"  json:"findingId"`
	CreatedAt  time.Time  `db:"created_at"  json:"createdAt"`
	UpdatedAt  *time.Time `db:"updated_at"  json:"updatedAt"`
}

// TagEvidenceMap reflects the structure of the database table 'tag_evidence_map'
type TagEvidenceMap struct {
	TagID      int64      `db:"tag_id"      json:"tagId"`
	EvidenceID int64      `db:"evidence_id" json:"evidenceId"`
	CreatedAt  time.Time  `db:"created_at"  json:"createdAt"`
	UpdatedAt  *time.Time `db:"updated_at"  json:"updatedAt"`
}

// Operation reflects the structure of the database table 'operations'
type Operation struct {
	ID        int64           `db:"id"         json:"id"`
	Slug      string          `db:"slug"       json:"slug"`
	Name      string          `db:"name"       json:"name"`
	Status    OperationStatus `db:"status"     json:"status"`
	CreatedAt time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time      `db:"updated_at" json:"updatedAt"`
}

type OperationStatus = int

const (
	OperationStatusPlanning OperationStatus = 0
	OperationStatusAcitve   OperationStatus = 1
	OperationStatusComplete OperationStatus = 2
)

// Tag reflects the structure of the database table 'tags'
type Tag struct {
	ID          int64      `db:"id"           json:"id"`
	OperationID int64      `db:"operation_id" json:"operationId"`
	Name        string     `db:"name"         json:"name"`
	ColorName   string     `db:"color_name"   json:"colorName"`
	CreatedAt   time.Time  `db:"created_at"   json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at"   json:"updatedAt"`
}

// User reflects the structure of the database table 'user'
type User struct {
	ID        int64      `db:"id"         json:"id"`
	Slug      string     `db:"slug"       json:"slug"`
	FirstName string     `db:"first_name" json:"firstName"`
	LastName  string     `db:"last_name"  json:"lastName"`
	Email     string     `db:"email"      json:"email"`
	Admin     bool       `db:"admin"      json:"admin"`
	Disabled  bool       `db:"disabled"   json:"disabled"`
	Headless  bool       `db:"headless"   json:"headless"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

// UserOperationPermission reflects the structure of the database table 'user_operation_permissions'
type UserOperationPermission struct {
	UserID      int64                `db:"user_id"`
	OperationID int64                `db:"operation_id"`
	Role        policy.OperationRole `db:"role"`
	CreatedAt   time.Time            `db:"created_at"`
	UpdatedAt   *time.Time           `db:"updated_at"`
}

// Query reflects the structure of the database table 'queries'
type Query struct {
	ID          int64      `db:"id"           json:"id"`
	OperationID int64      `db:"operation_id" json:"operationId"`
	Name        string     `db:"name"         json:"name"`
	Query       string     `db:"query"        json:"query"`
	Type        string     `db:"type"         json:"type"`
	CreatedAt   time.Time  `db:"created_at"   json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at"   json:"updatedAt"`
}

// AuthSchemeData reflects the structure of the database table 'auth_scheme_data'
type AuthSchemeData struct {
	ID                int64      `db:"id"`
	AuthScheme        string     `db:"auth_scheme"`
	UserKey           string     `db:"user_key"`
	UserID            int64      `db:"user_id"`
	EncryptedPassword []byte     `db:"encrypted_password"`
	MustResetPassword bool       `db:"must_reset_password"`
	LastLogin         *time.Time `db:"last_login"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
}

// Session reflects the structure of the database table 'sessions'
type Session struct {
	ID          int64      `db:"id"`
	UserID      int64      `db:"user_id"`
	SessionData []byte     `db:"session_data"`
	CreatedAt   time.Time  `db:"created_at"`
	ModifiedAt  *time.Time `db:"modified_at"`
	ExpiresAt   time.Time  `db:"expires_at"`
}

// ExportQueueItem reflects the structure of the database table 'exports_queue'
type ExportQueueItem struct {
	ID          int64        `db:"id"`
	OperationID int64        `db:"operation_id"`
	UserID      int64        `db:"user_id"` //Requesting user
	ExportName  string       `db:"export_name"`
	Status      ExportStatus `db:"status"`
	ErrorNotes  string       `db:"notes"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   *time.Time   `db:"updated_at"`
}

// ExportStatus reflects the "Status" field for ExportQueueItem
type ExportStatus = int

const (
	// ExportStatusPending is the default status that represents an export that has been queued, but not executed
	ExportStatusPending ExportStatus = 0
	// ExportStatusInProgress represents an export that has been started, but not completed
	ExportStatusInProgress ExportStatus = 1
	// ExportStatusComplete represents an export that has completed
	ExportStatusComplete ExportStatus = 2
	// ExportStatusError represents an export that could not complete due to an error (check the ErrorNotes)
	ExportStatusError ExportStatus = 3
	// ExportStatusCancelled represents an export that was cancelled by the user.
	ExportStatusCancelled ExportStatus = 4
)
