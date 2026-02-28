package remux

import (
	"net/http"

	"github.com/ashirt-ops/ashirt-server/internal/server/dissectors"
	"github.com/go-chi/chi/v5"
)

func generateUrlParamMap(r *http.Request) map[string]string {
	urlParamMap := map[string]string{}
	cxt := chi.RouteContext(r.Context())
	if cxt != nil {
		urlParams := cxt.URLParams
		for index, value := range urlParams.Keys {
			urlParamMap[value] = urlParams.Values[index]
		}
	}
	return urlParamMap
}

// DissectJSONRequest is a gorilla.mux focused rewrap of dissectors.DissectJSONRequest
func DissectJSONRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectJSONRequest(r, generateUrlParamMap(r))
}

// DissectFormRequest is a gorilla.mux focused rewrap of dissectors.DissectFormRequest
func DissectFormRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectFormRequest(r, generateUrlParamMap(r))
}

// DissectNoBodyRequest is a gorilla.mux focused rewrap of dissectors.DissectPlainRequest
func DissectNoBodyRequest(r *http.Request) dissectors.DissectedRequest {
	return dissectors.DissectPlainRequest(r, generateUrlParamMap(r))
}
