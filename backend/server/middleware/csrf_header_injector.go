// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func InjectCSRFTokenHeader() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-CSRF-Token", csrf.Token(r))
			next.ServeHTTP(w, r)
		})
	}
}
