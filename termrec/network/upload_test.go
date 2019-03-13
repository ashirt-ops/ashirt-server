package network_test

import (
	"bytes"
	"testing"

	"github.com/theparanoids/ashirt/termrec/network"
	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	var written []byte
	makeServer(Route{"POST", "/api/operations/777/evidence", newRequestRecorder(201, "", &written)})
	network.SetBaseURL("http://localhost" + testPort)

	uploadInput := network.UploadInput{
		OperationID: 777,
		Description: "abcd",
		Filename:    "dolphin",
		Content:     bytes.NewReader([]byte("abc123")),
	}

	err := network.UploadToAshirt(uploadInput)

	require.Nil(t, err)
}

func TestUploadFailedWithJSONError(t *testing.T) {
	var written []byte
	makeServer(Route{"POST", "/api/operations/778/evidence", newRequestRecorder(402, `{"error": "oops"}`, &written)})
	network.SetBaseURL("http://localhost" + testPort)

	uploadInput := network.UploadInput{
		OperationID: 778,
		Description: "abcd",
		Filename:    "dolphin",
		Content:     bytes.NewReader([]byte("abc123")),
	}

	err := network.UploadToAshirt(uploadInput)
	require.Error(t, err)
}

func TestUploadFailedWithUnknownJSON(t *testing.T) {
	var written []byte
	makeServer(Route{"POST", "/api/operations/776/evidence", newRequestRecorder(402, `{"something": "value"}`, &written)})
	network.SetBaseURL("http://localhost" + testPort)

	uploadInput := network.UploadInput{
		OperationID: 776,
		Description: "abcd",
		Filename:    "dolphin",
		Content:     bytes.NewReader([]byte("abc123")),
	}

	err := network.UploadToAshirt(uploadInput)
	require.Error(t, err)
}
