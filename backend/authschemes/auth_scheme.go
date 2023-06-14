// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package authschemes

import (
	"github.com/go-chi/chi/v5"
)

// AuthScheme provides a small interface into interacting with the AShirt backend authentication.
// The interface consists of two methods:
//
// Name() string: This method shall return a string that identifies the authentication scheme
// being used. It shall be distinct from any other authentication system being used within
// this project.
//
// FriendlyName() string: This method shall return a friendly version of the authentication that
// endusers will understand. It should, but is not strictly required, that the value be different
// from any other scheme. Likewise, it should be a "friendlier" version of Name(), though it need
// not be.
//
// BindRoutes(router, authBridge): BindRoutes exposes a _namespaced_ router that the authentication
// system can use to register custom endpoints. Each router is prefixed with /auth/{name} (as
// determined by the Name() method)
type AuthScheme interface {
	BindRoutes(chi.Router, AShirtAuthBridge)
	Name() string
	FriendlyName() string
	Flags() []string

	// Type provides a way to identify how a scheme works apart from its name. Currently this has two
	// "categories". First is "oidc", which is used for any generic OIDC provider. Second is the name
	// of the method (e.g. "local"), which is used when there's no real alternative to speak of.
	Type() string
}
