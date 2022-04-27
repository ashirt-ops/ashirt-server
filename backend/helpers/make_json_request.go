package helpers

import (
	"io"
	"net/http"
)

var client = &http.Client{}

type ModifyReqFunc = func(req *http.Request) error

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

func AddHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}

func NoMod(req *http.Request) error {
	return nil
}
