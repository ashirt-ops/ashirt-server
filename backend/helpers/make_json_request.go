package helpers

import (
	"io"
	"net/http"
)

var client = &http.Client{}

type ModifyReqFunc = func(req *http.Request) error

// MakeJSONRequest makes a request with the content-type application/json, and an optional body
func MakeJSONRequest(method, url string, body io.Reader, updateRequest ModifyReqFunc) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if err = updateRequest(req); err != nil {
		return nil, err
	}

	return client.Do(req)
}

// AddHeaders adds headers to a request pre-flight
func AddHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}

// NoMod is a canned value that can be used for MakeJSONRequest's updateRequest parameter
func NoMod(req *http.Request) error {
	return nil
}
