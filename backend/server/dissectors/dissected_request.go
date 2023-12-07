// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package dissectors

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/ashirt-ops/ashirt-server/backend"
)

// DissectedRequest stores a parsed JSON body, the map of URL substituted values, and query string
// values, as well as if any errors were encountered during processing.
//
// Usage Example:
//
// // Chi
//
//	func(w http.ResponseWriter, r *http.Request) {
//	  parsedRequest := DissectJSONRequest(r, GenerateUrlParamMap(r))
//	  input := service.RepeatWordInput{
//			SomeString: parsedRequest.FromURL("someString").Required(true).AsString()
//	     Times: parsedRequest.FromQuery("times").OrDefault(2).AsInt64()
//	  }
//	  if parsedRequest.Error != nil {
//	     // process error
//	  }
//	  service.RepeatWord( input )
//	}
type DissectedRequest struct {
	Request     *http.Request
	bodyValues  map[string]interface{}
	urlValues   map[string]string
	queryValues url.Values
	Error       error
}

func makeDefaultRequest(r *http.Request, urlParameters map[string]string) DissectedRequest {
	return DissectedRequest{
		Request:     r,
		bodyValues:  make(map[string]interface{}),
		urlValues:   urlParameters,
		queryValues: r.URL.Query(),
		Error:       nil,
	}
}

// DissectJSONRequest retrieves query string values from the provided request and parses a JSON
// body, if present. Once this data is parsed, values can be retrieved by using FromQuery, FromURL,
// or FromBody
//
// If an error is encountered while trying to parse the json, it will be stored in the
// DissectedRequest.Error field
func DissectJSONRequest(r *http.Request, urlParameters map[string]string) DissectedRequest {
	rtn := makeDefaultRequest(r, urlParameters)

	if r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&rtn.bodyValues); err != nil {
			rtn.Error = backend.BadInputErr(err, "Invalid JSON body")
		}
	}

	return rtn
}

// DissectPlainRequest retrieves query string values from the provided request.
// Once this data is parsed, values can be retrieved by using FromQuery, or FromURL
//
// If an error is encountered while trying to parse the json, it will be stored in the
// DissectedRequest.Error field
func DissectPlainRequest(r *http.Request, urlParameters map[string]string) DissectedRequest {
	return makeDefaultRequest(r, urlParameters)
}

// DissectFormRequest retrieves query string values from the provided request and parses a multipart
// form body, if present. Once this data is parsed, values can be retrieved by using FromQuery, FromURL,
// FromBody or FromFile. In the case of files, the particular keys need to be "registered"
// by passing values into the fileKeys parameter. Then they can be recalled by passing the same key
//
// If an error is encountered while trying to parse the json, it will be stored in the
// DissectedRequest.Error field
func DissectFormRequest(r *http.Request, urlParameters map[string]string) DissectedRequest {
	rtn := makeDefaultRequest(r, urlParameters)

	if r.Body != http.NoBody {
		if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
			rtn.Error = backend.BadInputErr(err, "Unable to parse Form body")
		}
		//need to unwrap the post form to coerce it into a consistent format
		for key, value := range r.PostForm {
			rtn.bodyValues[key] = value
		}
	}

	return rtn
}

// FromBody attempts to retrieve a field from the parsed body.
// Returns a Coercable, which can then be transformed into the desired type.
// If an error has already been encountered, the resulting action is a no-op
func (m *DissectedRequest) FromBody(key string) *Coercable {
	if m.Error != nil {
		return makeErrCoercable(m)
	}
	value, found := m.bodyValues[key]
	return makeCoercable(m, key, found, value)
}

// FromURL attempts to retrieve a field from the provided url parameters
// (e.g. in POST /operation/{operation_id}, operation_id would be the URL parameter).
// Returns Coercable, which can then be transformed into the desired type.
// If an error has already been encountered, the resulting action is a no-op
func (m *DissectedRequest) FromURL(key string) *Coercable {
	if m.Error != nil {
		return makeErrCoercable(m)
	}
	value, found := m.urlValues[key]
	return makeCoercable(m, key, found, value)
}

// FromQuery attempts to retrieve a field from the query string.
// Returns a Coercable, which can then be transformed into the desired type.
// If an error has already been encountered, the resulting action is a no-op
//
// Note: All values retrieved from the query string are string slices. Any method
// that would return a singular item instead returns the first element instead. Likewise,
// any method that returns a non-string will instead convert the field appropriately
func (m *DissectedRequest) FromQuery(key string) *Coercable {
	if m.Error != nil {
		return makeErrCoercable(m)
	}
	value, found := m.queryValues[key]
	return makeCoercable(m, key, found, value)
}

// FromFile retrieves the named file from a multipart form. The returned value is a
// *UploadedFile which can be used to retrieve the file or file header.
func (m *DissectedRequest) FromFile(key string) multipart.File {
	file, _, _ := m.Request.FormFile(key)
	return file
}
