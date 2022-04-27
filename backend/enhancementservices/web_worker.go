package enhancementservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
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
	w.Config = webConfig
	w.EvidenceID = evidenceID
	return nil
}

func (w *webConfigV1Worker) Test() (string, bool, error) {
	body := []byte(`{"type": "test"}`)
	resp, err := helpers.MakeJSONRequest("POST", w.Config.URL, bytes.NewReader(body), func(req *http.Request) error {
		helpers.AddHeaders(req, w.Config.Headers)
		return nil
	})
	if err != nil {
		return "Unable to verify worker status", false, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return "Service is functional", true, nil
	} else {
		var parsedData webTestResp
		if err := json.NewDecoder(resp.Body).Decode(&parsedData); err != nil {
			return "Unable to parse response", false, err
		}
		if parsedData.Status == "ok" {
			return "Service is functional", true, nil
		}
		if parsedData.Status == "error" {
			if parsedData.Message != nil {
				return *parsedData.Message, false, nil
			}
			return "Service is reporting an error", false, nil
		}
	}

	return "Service did not reply with a supported status", false, nil
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
		dbModel.Status = helpers.StringPtr("Unaccepted")
		dbModel.CanProcess = helpers.BoolPtr(false)
		dbModel.LastRunMessage = message
	}
	recordError := func(message *string) {
		dbModel.Status = helpers.StringPtr("Error")
		dbModel.LastRunMessage = message
	}
	recordDeferral := func() {
		dbModel.Status = helpers.StringPtr("Queued")
		dbModel.CanProcess = helpers.BoolPtr(true)
	}
	recordProcessed := func(content string) {
		dbModel.Status = helpers.StringPtr("Completed")
		dbModel.CanProcess = helpers.BoolPtr(true)
		dbModel.Body = content
	}

	var parsedData webProcessResp
	if err := json.NewDecoder(resp.Body).Decode(&parsedData); err != nil {
		recordError(helpers.StringPtr("Unable to parse response"))
		return
	}

	switch resp.StatusCode {
	case http.StatusOK: // 200
		switch parsedData.Action {
		case "procssed":
			if parsedData.Content != nil {
				recordProcessed(*parsedData.Content)
			} else {
				recordError(helpers.StringPtr("Content was not delivered for successful run"))
			}
		case "rejected":
			recordRejection(parsedData.Content)
		case "error":
			recordError(parsedData.Content)
		case "deferred":
			recordDeferral()
		default:
			recordError(helpers.SprintfPtr("Unexpceted response format (%v)", parsedData.Action))
		}
	case http.StatusAccepted:
		recordDeferral()
	case http.StatusNotAcceptable:
		recordRejection(nil)
	case http.StatusInternalServerError:
		recordError(nil)
	default:
		recordError(helpers.SprintfPtr("Unexpceted response status code (%v)", resp.StatusCode))
	}
}
