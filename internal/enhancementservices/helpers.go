package enhancementservices

import (
	"net/http"

	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/servicetypes/evidencemetadata"
)

func handleProcessResponse(dbModel *models.EvidenceMetadata, statusCode int, parsedData ProcessResponse) {
	switch statusCode {
	case http.StatusOK: // 200
		switch parsedData.Action {
		case "processed":
			if parsedData.Content != nil {
				recordProcessed(dbModel, *parsedData.Content)
			} else {
				recordError(dbModel, helpers.Ptr("Content was not delivered for successful run"))
			}
		case "rejected":
			recordRejection(dbModel, parsedData.Content)
		case "error":
			recordError(dbModel, parsedData.Content)
		case "deferred":
			recordDeferral(dbModel)
		default:
			recordError(dbModel, helpers.SprintfPtr("Unexpected response format (%v)", parsedData.Action))
		}
	case http.StatusAccepted:
		recordDeferral(dbModel)
	case http.StatusNotAcceptable:
		recordRejection(dbModel, nil)
	case http.StatusInternalServerError:
		recordError(dbModel, nil)
	default:
		recordError(dbModel, helpers.SprintfPtr("Unexpected response status code (%v)", statusCode))
	}
}

func recordRejection(dbModel *models.EvidenceMetadata, message *string) {
	dbModel.Status = evidencemetadata.StatusCompleted.Ptr()
	dbModel.CanProcess = helpers.Ptr(false)
	dbModel.LastRunMessage = message
}

func recordError(dbModel *models.EvidenceMetadata, message *string) {
	dbModel.Status = evidencemetadata.StatusError.Ptr()
	dbModel.LastRunMessage = message
}

func recordDeferral(dbModel *models.EvidenceMetadata) {
	dbModel.Status = evidencemetadata.StatusQueued.Ptr()
	dbModel.CanProcess = helpers.Ptr(true)
}

func recordProcessed(dbModel *models.EvidenceMetadata, content string) {
	dbModel.Status = evidencemetadata.StatusCompleted.Ptr()
	dbModel.CanProcess = helpers.Ptr(true)
	dbModel.Body = content
}
