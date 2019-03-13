// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package remux

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route rewraps gorilla.mux's Handle/Methods to provide a better at-a-glance reading of route definitions
func Route(r *mux.Router, method string, path string, handler http.Handler) {
	r.Handle(path, handler).Methods(method)
}
