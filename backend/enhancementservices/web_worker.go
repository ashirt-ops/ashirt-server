// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
)

type webConfigV1Worker struct {
	Config     WebConfigV1
	WorkerName string
	// makeRequestFn provides an alternative function to make a JSON based request. Should typically be nil,
	// except when unit testing
	makeRequestFn RequestFn
}

type WebConfigV1 struct {
	BasicServiceWorkerConfig
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func (w *webConfigV1Worker) Build(workerName string, workerConfig []byte) error {
	var webConfig WebConfigV1
	if err := json.Unmarshal([]byte(workerConfig), &webConfig); err != nil {
		return backend.WrapError("worker configuration is unparsable", err)
	}
	w.WorkerName = workerName
	w.Config = webConfig
	return nil
}

func (w *webConfigV1Worker) Test() ServiceTestResult {
	body := []byte(`{"type": "test"}`)
	resp, err := w.makeJSONRequest("POST", w.Config.URL, bytes.NewReader(body), func(req *http.Request) error {
		helpers.AddHeaders(req, w.Config.Headers)
		return nil
	})
	if err != nil {
		return ErrorTestResultWithMessage(err, "Unable to verify worker status")
	}

	if resp.StatusCode == http.StatusNoContent {
		return TestResultSuccess("Service is functional")
	} else {
		var parsedData TestResp
		if err := json.NewDecoder(resp.Body).Decode(&parsedData); err != nil {
			return ErrorTestResultWithMessage(err, "Unable to parse response")
		}
		if parsedData.Status == "ok" {
			return TestResultSuccess("Service is functional")
		}
		if parsedData.Status == "error" {
			if parsedData.Message != nil {
				return ErrorTestResultWithMessage(nil, *parsedData.Message)
			}
			return ErrorTestResultWithMessage(nil, "Service reported an error")
		}
	}

	return ErrorTestResultWithMessage(nil, "Service did not reply with a supported status")
}

func (w *webConfigV1Worker) Process(evidenceID int64, payload *Payload) (*models.EvidenceMetadata, error) {
	body, err := json.Marshal(*payload)
	if err != nil {
		return nil, backend.WrapError("unable to construct body", err)
	}

	resp, err := w.makeJSONRequest("POST", w.Config.URL, bytes.NewReader(body), func(req *http.Request) error {
		helpers.AddHeaders(req, w.Config.Headers)
		return nil
	})

	if err != nil {
		return nil, err
	}

	model := models.EvidenceMetadata{
		Source:     w.WorkerName,
		EvidenceID: evidenceID,
	}
	handleWebResponse(&model, resp)

	return &model, nil
}

func handleWebResponse(dbModel *models.EvidenceMetadata, resp *http.Response) {
	var parsedData ProcessResponse

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		recordError(dbModel, helpers.Ptr("Unable to read response"))
		return
	}
	if len(bytes) > 0 {
		if err := json.Unmarshal(bytes, &parsedData); err != nil {
			recordError(dbModel, helpers.Ptr("Unable to parse response"))
			return
		}
	}

	handleProcessResponse(dbModel, resp.StatusCode, parsedData)
}

func BuildTestWebWorker() webConfigV1Worker {
	return webConfigV1Worker{
		WorkerName: "magic",
		Config: WebConfigV1{
			URL:     "http://localhost/failifcalled",
			Headers: map[string]string{},
			BasicServiceWorkerConfig: BasicServiceWorkerConfig{
				Type:    "aws",
				Version: 1,
			},
		},
	}
}

func (l webConfigV1Worker) makeJSONRequest(method, url string, body io.Reader, updateRequest helpers.ModifyReqFunc) (*http.Response, error) {
	if l.makeRequestFn != nil {
		return l.makeRequestFn(method, url, body, updateRequest)
	}
	return helpers.MakeJSONRequest(method, url, body, updateRequest)
}

func (l *webConfigV1Worker) SetWebRequestFunction(fn RequestFn) {
	l.makeRequestFn = fn
}
