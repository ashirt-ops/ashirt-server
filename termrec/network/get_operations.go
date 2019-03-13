package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const errCannotConnectMsg = "Unable to connect to the server"

// GetOperations retrieves all of the operations that are exposed to backend tools (api routes)
// This should be replaced with a login and a web query once security is in place
func GetOperations() ([]Operation, error) {
	var ops []Operation
	req, err := http.NewRequest("GET", apiURL+"/operations", http.NoBody)

	if err != nil {
		return ops, errors.Wrap(err, errCannotConnectMsg)
	}

	err = addAuthentication(req)
	if err != nil {
		return ops, errors.Wrap(err, errCannotConnectMsg)
	}

	resp, err := client.Do(req)
	if err != nil {
		return ops, errors.Wrap(err, errCannotConnectMsg)
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized:
		return ops, errors.New("Unable to authenticate with server. Please check credentials")
	case resp.StatusCode == http.StatusInternalServerError:
		return ops, errors.New("Server encountered an error")
	case resp.StatusCode != http.StatusOK:
		return ops, errors.Wrap(err, errCannotConnectMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ops, errors.Wrap(err, "Unable to read response")
	}

	if err := json.Unmarshal(body, &ops); err != nil {
		return ops, errors.Wrap(err, "Unable to parse response")
	}

	return ops, nil
}
