package network

import (
	"net/http"
	"time"

	"github.com/theparanoids/ashirt/shared"
)

var client = &http.Client{}

var apiURL string
var accessKey string
var secretKey []byte

// SetBaseURL Sets the url to use as a base for all service contact
// Note: this function only requires the url to reach the frontend service.
// routes will be deduced from that.
func SetBaseURL(url string) {
	apiURL = url + "/api"
}

// BaseURLSet is a small check to verify that some value exists for the BaseURL
func BaseURLSet() bool {
	return apiURL != ""
}

// SetAccessKey sets the common access key for all API actions
func SetAccessKey(key string) {
	accessKey = key
}

// SetSecretKey sets the common secret key for all API actions
func SetSecretKey(key []byte) {
	secretKey = key
}

// addAuthentication adds Date and Authentication headers to the provided request
// returns an error if building an appropriate authentication value fails, nil otherwise
// Note: This should be called immediately before sending a request.
func addAuthentication(req *http.Request) error {
	req.Header.Set("Date", time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123))
	authorization, err := shared.BuildClientRequestAuthorization(req, accessKey, secretKey)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", authorization)
	return nil
}
