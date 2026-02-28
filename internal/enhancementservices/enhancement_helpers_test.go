package enhancementservices_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	// aliasing as this to shorten lines / aid in reading
	this "github.com/ashirt-ops/ashirt-server/internal/enhancementservices"
)

type transformFn = func(t *testing.T, data interface{}) []byte
type responseFn = func(t *testing.T, w *httptest.ResponseRecorder, content string, tf transformFn)

func testSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, tf transformFn) {
	w.WriteHeader(http.StatusOK)
	resp := this.TestResp{
		Status: "ok",
	}
	w.Write(tf(t, resp))
}

func testErrorResponse(t *testing.T, w *httptest.ResponseRecorder, message string, tf transformFn) {
	w.WriteHeader(http.StatusOK)
	resp := this.TestResp{
		Status:  "error",
		Message: &message,
	}
	w.Write(tf(t, resp))
}

func processSuccessReponse(t *testing.T, w *httptest.ResponseRecorder, content string, tf transformFn) {
	w.WriteHeader(http.StatusOK)
	resp := awsProcessResp{
		Action:  "processed",
		Content: &content,
	}
	w.Write(tf(t, resp))
}

func processErrorResponse_NoContent(t *testing.T, w *httptest.ResponseRecorder, _ string, tf transformFn) {
	w.WriteHeader(http.StatusOK)
	resp := awsProcessResp{
		Action:  "processed",
		Content: nil,
	}
	w.Write(tf(t, resp))
}

func processErrorResponse_WithMessage(t *testing.T, w *httptest.ResponseRecorder, content string, tf transformFn) {
	w.WriteHeader(http.StatusOK)
	resp := awsProcessResp{
		Action:  "error",
		Content: &content,
	}
	w.Write(tf(t, resp))
}

func processErrorResponse_StatusCode(t *testing.T, w *httptest.ResponseRecorder, _ string, _ transformFn) {
	w.WriteHeader(http.StatusInternalServerError)
}

func processDeferalResponse_StatusCode(t *testing.T, w *httptest.ResponseRecorder, _ string, _ transformFn) {
	w.WriteHeader(http.StatusAccepted)
}

func processDeferalReponse_Action(t *testing.T, w *httptest.ResponseRecorder, content string, tf transformFn) {
	w.WriteHeader(http.StatusOK)
	resp := awsProcessResp{
		Action: "deferred",
	}
	w.Write(tf(t, resp))
}

func processRejectedResponse_StatusCode(t *testing.T, w *httptest.ResponseRecorder, _ string, _ transformFn) {
	w.WriteHeader(http.StatusNotAcceptable)
}

func processRejectedReponse_Action(t *testing.T, w *httptest.ResponseRecorder, content string, tf transformFn) {
	w.WriteHeader(http.StatusOK)
	resp := awsProcessResp{
		Action:  "rejected",
		Content: &content,
	}
	w.Write(tf(t, resp))
}

func verifyTestBody(t *testing.T, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	container := make(map[string]interface{})
	json.Unmarshal(data, &container)
	typeVal, ok := container["type"]
	require.True(t, ok)
	require.Equal(t, "test", typeVal)
}

func verifyProcessBody(t *testing.T, req *http.Request, expectedPayload this.NewEvidencePayload) {
	var processBody this.NewEvidencePayload
	data, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &processBody)
	require.NoError(t, err)
	require.Equal(t, expectedPayload, processBody)
}

func wrapInAwsResponse(t *testing.T, data interface{}) []byte {
	body, err := json.Marshal(data)
	require.NoError(t, err)
	resp := this.LambdaResponse{
		Body:       string(body),
		StatusCode: 200,
	}
	respondWith, _ := json.Marshal(resp)
	return respondWith
}

func noWrap(t *testing.T, data interface{}) []byte {
	body, err := json.Marshal(data)
	require.NoError(t, err)
	return body
}

func buildLambdaClientWithResponse(respFn func(req *http.Request, err error) (*http.Response, error)) this.LambdaInvokableClient {
	return this.NewTestRIELambdaClient(
		makeMockRequestHandler(
			RequestMock{OnSendRequest: respFn},
		),
	)
}
