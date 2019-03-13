// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package remux

import (
	"net/http"

	"github.com/theparanoids/ashirt/backend/server/dissectors"
	"github.com/gorilla/mux"
)

// DissectJSONRequest is a gorilla.mux focused rewrap of dissectors.DissectJSONRequest
func DissectJSONRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectJSONRequest(r, mux.Vars(r))
}

// DissectFormRequest is a gorilla.mux focused rewrap of dissectors.DissectFormRequest
func DissectFormRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectFormRequest(r, mux.Vars(r))
}

// DissectNoBodyRequest is a gorilla.mux focused rewrap of dissectors.DissectPlainRequest
func DissectNoBodyRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectPlainRequest(r, mux.Vars(r))
}
