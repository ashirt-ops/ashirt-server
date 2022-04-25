package helpers

import (
	"io"
	"net/http"
	"time"

	"github.com/theparanoids/ashirt-server/signer"
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

func AddAShirtHMAC(req *http.Request, accessKey string, secretKey []byte) error {
	req.Header.Set("Date", time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123))
	authorization, err := signer.BuildClientRequestAuthorization(req, accessKey, secretKey)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", authorization)
	return nil
}

func AddHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}
