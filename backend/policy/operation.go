// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package policy

import "fmt"

type OperationRole string

const (
	OperationRoleAdmin OperationRole = "admin"
	OperationRoleWrite OperationRole = "write"
	OperationRoleRead  OperationRole = "read"
)

// Operation Policy
// Grants permissions based on operation roles
type Operation struct {
	UserID           int64
	IsHeadless       bool
	OperationRoleMap map[int64]OperationRole
}

func (o *Operation) String() string {
	return fmt.Sprintf("OperationPolicy(userID:%d, %v)", o.UserID, o.OperationRoleMap)
}

func (o *Operation) Check(permission Permission) bool {
	switch p := permission.(type) {
	case CanModifyUserOfOperation:
		return p.UserID != o.UserID && // A user cannot modify their own permissions (to prevent lockout)
			o.hasRole(p.OperationID, OperationRoleAdmin)

	case CanDeleteOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin)

	case CanModifyFindingsOfOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin, OperationRoleWrite) || o.IsHeadless
	case CanModifyEvidenceOfOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin, OperationRoleWrite) || o.IsHeadless
	case CanModifyOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin, OperationRoleWrite)
	case CanModifyQueriesOfOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin, OperationRoleWrite) || o.IsHeadless
	case CanModifyTagsOfOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin, OperationRoleWrite) || o.IsHeadless

	case CanListUsersOfOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin, OperationRoleWrite, OperationRoleRead) || o.IsHeadless
	case CanReadOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin, OperationRoleWrite, OperationRoleRead) || o.IsHeadless

	case CanListUserGroupsOfOperation:
		return o.hasRole(p.OperationID, OperationRoleAdmin) || o.IsHeadless
	}
	return false
}

func (o *Operation) hasRole(operationID int64, allowedRoles ...OperationRole) bool {
	role, ok := o.OperationRoleMap[operationID]
	if !ok {
		return false
	}

	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}
