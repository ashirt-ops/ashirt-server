package network_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type handler func(http.ResponseWriter, *http.Request)

const testPort = ":12345"

type Route struct {
	Method  string
	Path    string
	Handler handler
}

func newRoute(method, path string, h handler) Route {
	return Route{Method: method, Path: path, Handler: h}
}

func newCannedResponse(status int, resp string) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(resp))
	}
}

func newRequestRecorder(status int, resp string, body *[]byte) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		*body, _ = ioutil.ReadAll(r.Body)
		w.WriteHeader(status)
		w.Write([]byte(resp))
	}
}

func makeServer(routes ...Route) {
	for _, r := range routes {
		http.HandleFunc(r.Path, r.Handler)
	}

	go func() {
		log.Fatal(http.ListenAndServe(testPort, nil))
	}()
}

func TestNetworkTestHelper_Gets(t *testing.T) {
	t.Skip("skipping network tests")
	makeServer(Route{"GET", "/hi", newCannedResponse(200, "hello!")})
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost"+testPort+"/hi", http.NoBody)
	assert.Nil(t, err)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "hello!", string(body))
}

func TestNetworkTestHelper_Posts(t *testing.T) {
	t.Skip("skipping network tests")
	msg := []byte("ABC123")
	var written []byte
	makeServer(Route{"GET", "/bye", newRequestRecorder(201, "late, yo", &written)})
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost"+testPort+"/bye", bytes.NewReader(msg))
	assert.Nil(t, err)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "late, yo", string(body))
	assert.Equal(t, msg, written)
}
