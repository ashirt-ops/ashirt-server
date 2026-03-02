package remux

import (
	"net/http"
)

// Route registers an HTTP handler for the given method and path on mux.
// It combines method and path into the stdlib pattern format "METHOD /path".
func Route(mux *http.ServeMux, method string, path string, handler http.Handler) {
	mux.Handle(method+" "+path, handler)
}
