// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package policy

import (
	"fmt"
	"reflect"

	recoveryConsts "github.com/theparanoids/ashirt/backend/authschemes/recoveryauth/constants"
)

// Policy is a simple interface into interacting with Permission structs.
//
// Check verifies, given a Permission, whether the requested action is valid/authorized
type Policy interface {
	Check(Permission) bool
	String() string
}

// Require institutes a policy check. Require will implicitly call policy.Check on the provided
// policy for each requiredPermission passed. If any of the required permissions fail the Check call,
// then this function will return an error. A response of (nil) means that all checks passed.
func Require(policy Policy, requiredPermissions ...Permission) error {
	for _, permission := range requiredPermissions {
		if !policy.Check(permission) {
			return fmt.Errorf("%s failed permission check %s%v", policy.String(), reflect.TypeOf(permission), permission)
		}
	}
	return nil
}

// Deny Policy
// Rejects all permissions
type Deny struct{}

func (*Deny) String() string { return "DenyPolicy" }

// Check returns false for every input, simulating a Never-Allow scenario
func (*Deny) Check(Permission) bool { return false }

// FullAccess Policy
// Allows all permissions
type FullAccess struct{}

func (*FullAccess) String() string { return "FullAccessPolicy" }

// Check returns true for every input, simulating an Always-Allow scenario
func (*FullAccess) Check(Permission) bool { return true }

// Union Policy
// Grants permission if either sub policy grants the permission
type Union struct {
	P1 Policy
	P2 Policy
}

func (u *Union) String() string { return fmt.Sprintf("%s|%s", u.P1.String(), u.P2.String()) }

// Check performs the underlying policy check for each policy, returning true if either
// policy.Check call returns true
func (u *Union) Check(permission Permission) bool {
	return u.P1.Check(permission) || u.P2.Check(permission)
}

// Authenticated Policy
// Grants permissions all authenticated users should have
type Authenticated struct {
	UserID       int64
	IsSuperAdmin bool
}

func NewAuthenticatedPolicy(userID int64, isSuperAdmin bool) *Authenticated {
	a := Authenticated{
		UserID:       userID,
		IsSuperAdmin: isSuperAdmin,
	}
	return &a
}

func (a *Authenticated) String() string {
	return fmt.Sprintf("AuthenticatedPolicy(userID:%d)", a.UserID)
}

// Check reviews the permission type for all authenticated users (true => valid ;; false => invalid)
func (a *Authenticated) Check(permission Permission) bool {
	switch target := permission.(type) {
	case CanCreateOperations:
		return true

	case AdminUsersOnly:
		return a.IsSuperAdmin

	case CanModifyAPIKeys:
		return target.UserID == a.UserID || a.IsSuperAdmin
	case CanListAPIKeys:
		return target.UserID == a.UserID || a.IsSuperAdmin

	case CanDeleteAuthScheme:
		return (target.UserID == a.UserID || a.IsSuperAdmin) && target.SchemeCode != recoveryConsts.Code
	case CanDeleteAuthForAllUsers:
		return a.IsSuperAdmin && target.SchemeCode != recoveryConsts.Code

	case CanModifyUser:
		return target.UserID == a.UserID || a.IsSuperAdmin
	case CanReadDetailedUser:
		return target.UserID == a.UserID || a.IsSuperAdmin
	case CanReadUser:
		return true
	}
	return false
}
