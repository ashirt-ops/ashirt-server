// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	// aliasing as this to shorten lines / aid in reading
	this "github.com/theparanoids/ashirt-server/backend/enhancementservices"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/servicetypes/evidencemetadata"
)

func TestAWSTest(t *testing.T) {
	worker := this.BuildTestLambdaWorker()

	setClient := func(writeResponse func(w *httptest.ResponseRecorder)) {
		this.SetTestLambdaClient(
			buildLambdaClientWithResponse(func(req *http.Request, err error) (*http.Response, error) {
				require.NoError(t, err)
				verifyTestBody(t, req)

				w := httptest.NewRecorder()
				writeResponse(w)
				return w.Result(), nil
			}),
		)
	}

	// verify test success
	setClient(func(w *httptest.ResponseRecorder) {
		testSuccessResponse(t, w, wrapInAwsResponse)
	})
	result := worker.Test()
	require.Equal(t, true, result.Live)
	require.NoError(t, result.Error)

	// verify test failure
	msg := "bummer"
	setClient(func(w *httptest.ResponseRecorder) {
		testErrorResponse(t, w, msg, wrapInAwsResponse)
	})
	result = worker.Test()
	require.Equal(t, false, result.Live)
	require.NoError(t, result.Error) // Error response aren't actually golang errors, so still expect noerror
	require.Contains(t, result.Message, msg)
}

func TestAWSProcess(t *testing.T) {
	worker := this.BuildTestLambdaWorker()

	payload := this.Payload{
		Type:          "process",
		EvidenceUUID:  "abc123",
		OperationSlug: "whatsit",
		ContentType:   "image",
	}
	content := "something cool"
	var eviID int64 = 123456

	buildProcessClient := func(fn responseFn) this.LambdaInvokableClient {
		return buildLambdaClientWithResponse(func(req *http.Request, err error) (*http.Response, error) {
			require.NoError(t, err)
			verifyProcessBody(t, req, payload)

			w := httptest.NewRecorder()
			fn(t, w, content, wrapInAwsResponse)
			return w.Result(), nil
		})
	}

	// verify success
	this.SetTestLambdaClient(buildProcessClient(processSuccessReponse))
	result, err := worker.Process(eviID, &payload)
	require.NoError(t, err)
	require.Equal(t, eviID, result.EvidenceID)
	require.True(t, *result.CanProcess)
	require.Equal(t, content, result.Body)
	require.Equal(t, evidencemetadata.StatusCompleted.Ptr(), result.Status)

	// verify Error Scenarios
	{
		verifyErrorScenario := func(result *models.EvidenceMetadata, err error) {
			require.NoError(t, err)
			require.Equal(t, eviID, result.EvidenceID)
			require.Equal(t, evidencemetadata.StatusError.Ptr(), result.Status)
		}

		// no-content failure
		this.SetTestLambdaClient(buildProcessClient(processErrorResponse_NoContent))
		result, err = worker.Process(eviID, &payload)
		verifyErrorScenario(result, err)
		require.NotNil(t, result.LastRunMessage)

		// with message
		this.SetTestLambdaClient(buildProcessClient(processErrorResponse_WithMessage))
		result, err = worker.Process(eviID, &payload)
		verifyErrorScenario(result, err)
		require.Equal(t, content, *result.LastRunMessage)

		// without message
		this.SetTestLambdaClient(buildProcessClient(processErrorResponse_StatusCode))
		result, err = worker.Process(eviID, &payload)
		verifyErrorScenario(result, err)
		require.Nil(t, result.LastRunMessage)
	}

	// verify deferals
	{
		verifyDefferalResult := func(result *models.EvidenceMetadata, err error) {
			require.NoError(t, err)
			require.Equal(t, eviID, result.EvidenceID)
			require.Equal(t, evidencemetadata.StatusQueued.Ptr(), result.Status)
			require.Equal(t, true, *result.CanProcess)
		}
		// status code version
		this.SetTestLambdaClient(buildProcessClient(processDeferalResponse_StatusCode))
		verifyDefferalResult(worker.Process(eviID, &payload))

		// action version
		this.SetTestLambdaClient(buildProcessClient(processDeferalReponse_Action))
		verifyDefferalResult(worker.Process(eviID, &payload))
	}

	// verify Rejections
	{
		verifyRejectionResult := func(result *models.EvidenceMetadata, err error) {
			require.NoError(t, err)
			require.Equal(t, eviID, result.EvidenceID)
			require.Equal(t, evidencemetadata.StatusCompleted.Ptr(), result.Status)
		}
		// status code version
		this.SetTestLambdaClient(buildProcessClient(processRejectedResponse_StatusCode))
		verifyRejectionResult(worker.Process(eviID, &payload))

		// action version
		this.SetTestLambdaClient(buildProcessClient(processRejectedReponse_Action))
		result, err := worker.Process(eviID, &payload)
		verifyRejectionResult(result, err)
		require.Equal(t, content, *result.LastRunMessage)
	}
}

type awsProcessResp struct {
	Action  string  `json:"action,omitempty"`  // Rejected | Deferred | Processed | Error
	Content *string `json:"content,omitempty"` // Error => reason, Processed => Result
}

type awsResponseContainer struct {
	Body string `json:"body"`
}
