// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package policy

// Permission represents a low level action that a policy can enforce.
// Permissions may optionally reference resources such as an operation or user
type Permission interface{}

type CanCreateOperations struct{}

// AdminUsersOnly addresses scenarios where the *only* restriction is that the
// calling user be an admin.
type AdminUsersOnly struct{}

type CanModifyAPIKeys struct{ UserID int64 }
type CanReadUser struct{ UserID int64 }
type CanReadDetailedUser struct{ UserID int64 }
type CanModifyUser struct{ UserID int64 }
type CanListAPIKeys struct{ UserID int64 }
type CanCheckTotp struct{ UserID int64 }
type CanDeleteTotp struct{ UserID int64 }

type CanDeleteAuthScheme struct {
	UserID     int64
	SchemeCode string
}
type CanDeleteAuthForAllUsers struct{ SchemeCode string }

// TODO TN set these up for user gruops
type CanListUsersOfOperation struct{ OperationID int64 }
type CanModifyFindingsOfOperation struct{ OperationID int64 }
type CanModifyEvidenceOfOperation struct{ OperationID int64 }
type CanModifyOperation struct{ OperationID int64 }
type CanModifyQueriesOfOperation struct{ OperationID int64 }
type CanModifyTagsOfOperation struct{ OperationID int64 }
type CanReadOperation struct{ OperationID int64 }
type CanDeleteOperation struct{ OperationID int64 }
type CanModifyUserOfOperation struct {
	OperationID int64
	UserID      int64
}
