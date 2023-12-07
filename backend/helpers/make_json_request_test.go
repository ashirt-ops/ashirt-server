package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ashirt-ops/ashirt-server/backend/helpers"
)

func TestAddHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost", http.NoBody)

	headers := map[string]string{
		"key":  "value",
		"key2": "value2",
	}
	helpers.AddHeaders(req, headers)

	for k, v := range headers {
		require.Equal(t, v, req.Header.Get(k))
	}
}
