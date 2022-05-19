// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/servicetypes/evidencemetadata"
)

type webConfigV1Worker struct {
	EvidenceID int64
	Config     WebConfigV1
	WorkerName string
}

type WebConfigV1 struct {
	BasicServiceWorkerConfig
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type webProcessResp struct {
	Action  string  `json:"action"`  // Rejected | Deferred | Processed | Error
	Content *string `json:"content"` // Rejected/Error => reason, Deferred => null/ignored, Processed => Result
}

type webTestResp struct {
	Status  string  `json:"status"`
	Message *string `json:"message"`
}

func (w *webConfigV1Worker) Build(workerName string, evidenceID int64, workerConfig []byte) error {
	var webConfig WebConfigV1
	if err := json.Unmarshal([]byte(workerConfig), &webConfig); err != nil {
		return err
	}
	w.WorkerName = workerName
	w.Config = webConfig
	w.EvidenceID = evidenceID
	return nil
}

func (w *webConfigV1Worker) Test() ServiceTestResult {
	body := []byte(`{"type": "test"}`)
	resp, err := helpers.MakeJSONRequest("POST", w.Config.URL, bytes.NewReader(body), func(req *http.Request) error {
		helpers.AddHeaders(req, w.Config.Headers)
		return nil
	})
	if err != nil {
		return ErrorTestResultWithMessage(err, "Unable to verify worker status")
	}

	if resp.StatusCode == http.StatusNoContent {
		return TestResultSuccess("Service is functional")
	} else {
		var parsedData webTestResp
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
			return ErrorTestResultWithMessage(nil, "Service is reporting an error")
		}
	}

	return ErrorTestResultWithMessage(nil, "Service did not reply with a supported status")
}

func (w *webConfigV1Worker) Process(payload *Payload) (*models.EvidenceMetadata, error) {
	body, err := json.Marshal(*payload)
	if err != nil {
		return nil, fmt.Errorf("unable to construct body")
	}

	resp, err := helpers.MakeJSONRequest("POST", w.Config.URL, bytes.NewReader(body), func(req *http.Request) error {
		helpers.AddHeaders(req, w.Config.Headers)
		return nil
	})

	if err != nil {
		return nil, err
	}

	model := models.EvidenceMetadata{
		Source:     w.WorkerName,
		EvidenceID: w.EvidenceID,
	}
	handleWebResponse(&model, resp)

	return &model, nil
}

func handleWebResponse(dbModel *models.EvidenceMetadata, resp *http.Response) {
	recordRejection := func(message *string) {
		dbModel.Status = evidencemetadata.StatusUnaccepted.Ptr()
		dbModel.CanProcess = helpers.Ptr(false)
		dbModel.LastRunMessage = message
	}
	recordError := func(message *string) {
		dbModel.Status = evidencemetadata.StatusError.Ptr()
		dbModel.LastRunMessage = message
	}
	recordDeferral := func() {
		dbModel.Status = evidencemetadata.StatusQueued.Ptr()
		dbModel.CanProcess = helpers.Ptr(true)
	}
	recordProcessed := func(content string) {
		dbModel.Status = evidencemetadata.StatusCompleted.Ptr()
		dbModel.CanProcess = helpers.Ptr(true)
		dbModel.Body = content
	}

	var parsedData webProcessResp
	if err := json.NewDecoder(resp.Body).Decode(&parsedData); err != nil {
		recordError(helpers.Ptr("Unable to parse response"))
		return
	}

	switch resp.StatusCode {
	case http.StatusOK: // 200
		switch parsedData.Action {
		case "processed":
			if parsedData.Content != nil {
				recordProcessed(*parsedData.Content)
			} else {
				recordError(helpers.Ptr("Content was not delivered for successful run"))
			}
		case "rejected":
			recordRejection(parsedData.Content)
		case "error":
			recordError(parsedData.Content)
		case "deferred":
			recordDeferral()
		default:
			recordError(helpers.SprintfPtr("Unexpected response format (%v)", parsedData.Action))
		}
	case http.StatusAccepted:
		recordDeferral()
	case http.StatusNotAcceptable:
		recordRejection(nil)
	case http.StatusInternalServerError:
		recordError(nil)
	default:
		recordError(helpers.SprintfPtr("Unexpected response status code (%v)", resp.StatusCode))
	}
}
