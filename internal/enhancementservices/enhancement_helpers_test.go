package enhancementservices_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	// aliasing as this to shorten lines / aid in reading
	this "github.com/ashirt-ops/ashirt-server/internal/enhancementservices"
	"github.com/ashirt-ops/ashirt-server/internal/helpers"
)

type responseFn = func(t *testing.T, w *httptest.ResponseRecorder, content string)

func writeJSON(t *testing.T, w *httptest.ResponseRecorder, data interface{}) {
	body, err := json.Marshal(data)
	require.NoError(t, err)
	w.Write(body)
}

func testSuccessResponse(t *testing.T, w *httptest.ResponseRecorder) {
	w.WriteHeader(http.StatusOK)
	writeJSON(t, w, this.TestResp{
		Status: "ok",
	})
}

func testErrorResponse(t *testing.T, w *httptest.ResponseRecorder, message string) {
	w.WriteHeader(http.StatusOK)
	writeJSON(t, w, this.TestResp{
		Status:  "error",
		Message: &message,
	})
}

func processSuccessReponse(t *testing.T, w *httptest.ResponseRecorder, content string) {
	w.WriteHeader(http.StatusOK)
	writeJSON(t, w, processResp{
		Action:  "processed",
		Content: &content,
	})
}

func processErrorResponse_NoContent(t *testing.T, w *httptest.ResponseRecorder, _ string) {
	w.WriteHeader(http.StatusOK)
	writeJSON(t, w, processResp{
		Action:  "processed",
		Content: nil,
	})
}

func processErrorResponse_WithMessage(t *testing.T, w *httptest.ResponseRecorder, content string) {
	w.WriteHeader(http.StatusOK)
	writeJSON(t, w, processResp{
		Action:  "error",
		Content: &content,
	})
}

func processErrorResponse_StatusCode(t *testing.T, w *httptest.ResponseRecorder, _ string) {
	w.WriteHeader(http.StatusInternalServerError)
}

func processDeferalResponse_StatusCode(t *testing.T, w *httptest.ResponseRecorder, _ string) {
	w.WriteHeader(http.StatusAccepted)
}

func processDeferalReponse_Action(t *testing.T, w *httptest.ResponseRecorder, content string) {
	w.WriteHeader(http.StatusOK)
	writeJSON(t, w, processResp{
		Action: "deferred",
	})
}

func processRejectedResponse_StatusCode(t *testing.T, w *httptest.ResponseRecorder, _ string) {
	w.WriteHeader(http.StatusNotAcceptable)
}

func processRejectedReponse_Action(t *testing.T, w *httptest.ResponseRecorder, content string) {
	w.WriteHeader(http.StatusOK)
	writeJSON(t, w, processResp{
		Action:  "rejected",
		Content: &content,
	})
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

type processResp struct {
	Action  string  `json:"action,omitempty"`  // Rejected | Deferred | Processed | Error
	Content *string `json:"content,omitempty"` // Error => reason, Processed => Result
}

func makeMockRequestHandler(mock RequestMock) this.RequestFn {
	return func(method, url string, body io.Reader, updateRequest helpers.ModifyReqFunc) (*http.Response, error) {
		content, err := io.ReadAll(body)
		clonedBody := bytes.NewReader(content)
		req := httptest.NewRequest(method, url, clonedBody)

		if mock.OnInvoked != nil {
			mock.OnInvoked(RequestData{
				Method:  method,
				URL:     url,
				Body:    content,
				Request: req,
				Error:   err,
			})
		}
		err = updateRequest(req)
		if mock.OnSendRequest != nil {
			return mock.OnSendRequest(req, err)
		}

		// default in case someone doesn't provide a RespondWith function
		w := httptest.NewRecorder()
		w.WriteHeader(http.StatusNoContent)
		return w.Result(), nil
	}
}

// opting for a struct here so On* functions can be omitted
type RequestMock struct {
	OnInvoked     func(RequestData)
	OnSendRequest func(*http.Request, error) (*http.Response, error)
}

type RequestData struct {
	Method  string
	URL     string
	Body    []byte
	Request *http.Request
	Error   error
}
