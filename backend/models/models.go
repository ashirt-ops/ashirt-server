// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package models

import (
	"time"

	"github.com/theparanoids/ashirt-server/backend/policy"
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
	ID            int64      `db:"id"`
	UUID          string     `db:"uuid"`
	OperationID   int64      `db:"operation_id"`
	ReadyToReport bool       `db:"ready_to_report"`
	TicketLink    *string    `db:"ticket_link"`
	CategoryID    *int64     `db:"category_id"`
	Title         string     `db:"title"`
	Description   string     `db:"description"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

// Evidence reflects the structure of the database table 'evidence'
type Evidence struct {
	ID            int64      `db:"id"`
	UUID          string     `db:"uuid"`
	OperationID   int64      `db:"operation_id"`
	OperatorID    int64      `db:"operator_id"`
	Description   string     `db:"description"`
	ContentType   string     `db:"content_type"`
	FullImageKey  string     `db:"full_image_key"`
	ThumbImageKey string     `db:"thumb_image_key"`
	OccurredAt    time.Time  `db:"occurred_at"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type EvidenceMetadata struct {
	ID             int64      `db:"id"`
	EvidenceID     int64      `db:"evidence_id"`
	Source         string     `db:"source"`
	Body           string     `db:"body"`
	Status         *string    `db:"status"`
	LastRunMessage *string    `db:"last_run_message"`
	CanProcess     *bool      `db:"can_process"`
	CreatedAt      time.Time  `db:"created_at"`
	WorkStartedAt  *time.Time `db:"work_started_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

// EvidenceFindingMap reflects the structure of the database table 'evidence_finding_map'
type EvidenceFindingMap struct {
	EvidenceID int64      `db:"evidence_id"`
	FindingID  int64      `db:"finding_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

// TagEvidenceMap reflects the structure of the database table 'tag_evidence_map'
type TagEvidenceMap struct {
	TagID      int64      `db:"tag_id"`
	EvidenceID int64      `db:"evidence_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

// Operation reflects the structure of the database table 'operations'
type Operation struct {
	ID        int64           `db:"id"`
	Slug      string          `db:"slug"`
	Name      string          `db:"name"`
	Status    OperationStatus `db:"status"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt *time.Time      `db:"updated_at"`
}

type OperationStatus = int

const (
	OperationStatusPlanning OperationStatus = 0
	OperationStatusAcitve   OperationStatus = 1
	OperationStatusComplete OperationStatus = 2
)

// Tag reflects the structure of the database table 'tags'
type Tag struct {
	ID          int64      `db:"id"`
	OperationID int64      `db:"operation_id"`
	Name        string     `db:"name"`
	ColorName   string     `db:"color_name"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// User reflects the structure of the database table 'user'
type User struct {
	ID        int64      `db:"id"`
	Slug      string     `db:"slug"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	Email     string     `db:"email"`
	Admin     bool       `db:"admin"`
	Disabled  bool       `db:"disabled"`
	Headless  bool       `db:"headless"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
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
	ID          int64      `db:"id"`
	OperationID int64      `db:"operation_id"`
	Name        string     `db:"name"`
	Query       string     `db:"query"`
	Type        string     `db:"type"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// AuthSchemeData reflects the structure of the database table 'auth_scheme_data'
type AuthSchemeData struct {
	ID                int64      `db:"id"`
	AuthScheme        string     `db:"auth_scheme"`
	AuthType          string     `db:"auth_type"`
	UserKey           string     `db:"user_key"`
	UserID            int64      `db:"user_id"`
	EncryptedPassword []byte     `db:"encrypted_password"`
	MustResetPassword bool       `db:"must_reset_password"`
	TOTPSecret        *string    `db:"totp_secret"`
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

// QueuedEmail reflects the structure of the database table 'email_queue'
type QueuedEmail struct {
	ID          int64      `db:"id"`
	ToEmail     string     `db:"to_email"`
	UserID      int64      `db:"user_id"`
	Template    string     `db:"template"`
	EmailStatus string     `db:"email_status"`
	ErrorCount  int64      `db:"error_count"`
	ErrorText   *string    `db:"error_text"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type FindingCategory struct {
	ID        int64      `db:"id"`
	Category  string     `db:"category"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type DefaultTag struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	ColorName string     `db:"color_name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type ServiceWorker struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Config    string     `db:"config"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
