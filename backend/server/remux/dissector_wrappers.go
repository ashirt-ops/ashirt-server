package remux

import (
	"net/http"

	"github.com/ashirt-ops/ashirt-server/backend/server/dissectors"
)

// DissectJSONRequest is a rewrap of dissectors.DissectJSONRequest
func DissectJSONRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectJSONRequest(r)
}

// DissectFormRequest is a rewrap of dissectors.DissectFormRequest
func DissectFormRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectFormRequest(r)
}

// DissectNoBodyRequest is a rewrap of dissectors.DissectPlainRequest
func DissectNoBodyRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectPlainRequest(r)
}
