// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package remux

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/logging"
)

// MediaHandler provides a generic handler for any content that _prefers_ a return value as raw data.
// In success situtations, the response will be returned as just a stream of bytes. In failure cases
// the output will instead be json, which conforms to the remainder of the project
func MediaHandler(handler func(*http.Request) (io.Reader, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data io.Reader
		var err error
		defer watcher(logging.ReqLogger(r.Context()), func(paniced bool) {
			if paniced {
				err = backend.PanicedError()
			}
			if err != nil {
				HandleError(w, r, err)
				return
			}

			io.Copy(w, data)
		})
		data, err = handler(r)
	})
}

// JSONHandler provides a generic handler for any request that prefers JSON responses. In all
// success scenarios, and most error scenarios, json is returned. The exception here is when
// this project cannot decode/Marshal a JSON message, in which case a plain 500 error with no content
// is returned.
func JSONHandler(handler func(*http.Request) (interface{}, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		var err error
		defer watcher(logging.ReqLogger(r.Context()), func(paniced bool) {
			if paniced {
				err = backend.PanicedError()
			}
			if err != nil {
				HandleError(w, r, err)
				return
			}

			status := 200
			if r.Method == "POST" {
				status = 201
			}
			writeJSONResponse(w, status, data)
		})
		data, err = handler(r)
	})
}

// HandleError will set the proper status code for the given error and return a json response
// body with a public reason listed
//
// Note: In general, users should prefer to use JSONHandler or MediaHandler. This function should
// only be used in instances where those handlers cannot be used (e.g. because of a redirect)
func HandleError(w http.ResponseWriter, r *http.Request, rootErr error) {
	var status int
	var publicReason string
	var loggedReason error

	switch err := rootErr.(type) {
	case *backend.HTTPError:
		status = err.HTTPStatus
		publicReason = err.PublicReason
		loggedReason = err.WrappedError

	case error:
		status = http.StatusInternalServerError
		publicReason = "An unknown error occurred"
		loggedReason = err
		logging.Log(r.Context(), "msg", "handling non-HTTPError", "stacktrace", formatStackTrace(retrace(20)))
	}

	logging.Log(r.Context(),
		"msg", "Error handling request",
		"error", loggedReason,
		"rootError", rootErr,
		"status", status,
		"url", r.URL,
	)

	writeJSONResponse(w, status, map[string]string{"error": publicReason})
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}
