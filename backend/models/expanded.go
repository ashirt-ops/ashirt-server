// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package models

import "time"

// UserWithAuthData represents a limited joining of users table with auth_scheme_data table
type UserWithAuthData struct {
	User
	AuthSchemeData []LimitedAuthSchemeData
}

// LimitedAuthSchemeData represents a partial AuthSchemeData model, exposing only the name of the scheme
type LimitedAuthSchemeData struct {
	AuthScheme string
	LastLogin  *time.Time
}
