// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

/*
	This file provides a rewrap of the relevant remux package. This is done to make migration easier,
	and to minimize the impact of this refactor
*/
package server

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theparanoids/ashirt/backend/server/dissectors"
	"github.com/theparanoids/ashirt/backend/server/remux"
)

func route(r *mux.Router, method string, path string, handler http.Handler) {
	remux.Route(r, method, path, handler)
}

func dissectJSONRequest(r *http.Request) dissectors.DissectedRequest {
	return remux.DissectJSONRequest(r)
}

func dissectFormRequest(r *http.Request) dissectors.DissectedRequest {
	return remux.DissectFormRequest(r)
}

func dissectNoBodyRequest(r *http.Request) dissectors.DissectedRequest {
	return remux.DissectNoBodyRequest(r)
}

func mediaHandler(handler func(*http.Request) (io.Reader, error)) http.Handler {
	return remux.MediaHandler(handler)
}

func jsonHandler(handler func(*http.Request) (interface{}, error)) http.Handler {
	return remux.JSONHandler(handler)
}
