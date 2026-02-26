package remux

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route rewraps chi's Handle/Methods to provide a better at-a-glance reading of route definitions
func Route(r chi.Router, method string, path string, handler http.Handler) {
	r.Method(method, path, handler)
}
