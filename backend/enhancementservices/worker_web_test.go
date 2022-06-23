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
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/servicetypes/evidencemetadata"
)

func TestWebTest(t *testing.T) {
	worker := helpers.Ptr(this.BuildTestWebWorker())

	setClient := func(writeResponse func(w *httptest.ResponseRecorder)) {
		worker.SetWebRequestFunction(makeMockRequestHandler(RequestMock{
			OnSendRequest: func(req *http.Request, err error) (*http.Response, error) {
				require.NoError(t, err)
				verifyTestBody(t, req)

				w := httptest.NewRecorder()
				writeResponse(w)
				return w.Result(), nil
			},
		}))
	}

	// verify test success
	setClient(func(w *httptest.ResponseRecorder) {
		testSuccessResponse(t, w, noWrap)
	})
	result := worker.Test()
	require.Equal(t, true, result.Live)
	require.NoError(t, result.Error)

	// verify test failure
	msg := "bummer"
	setClient(func(w *httptest.ResponseRecorder) {
		testErrorResponse(t, w, msg, noWrap)
	})
	result = worker.Test()
	require.Equal(t, false, result.Live)
	require.NoError(t, result.Error) // Error response aren't actually golang errors, so still expect noerror
	require.Contains(t, result.Message, msg)
}

func TestWebProcessMetadata(t *testing.T) {
	worker := helpers.Ptr(this.BuildTestWebWorker())

	payload := this.NewEvidencePayload{
		Type:          "evidence_created",
		EvidenceUUID:  "abc123",
		OperationSlug: "whatsit",
		ContentType:   "image",
	}
	content := "something cool"
	var eviID int64 = 123456

	setClient := func(fn responseFn) {
		worker.SetWebRequestFunction(makeMockRequestHandler(RequestMock{
			OnSendRequest: func(req *http.Request, err error) (*http.Response, error) {
				require.NoError(t, err)
				verifyProcessBody(t, req, payload)

				w := httptest.NewRecorder()
				fn(t, w, content, noWrap)

				return w.Result(), nil
			},
		}))
	}

	// verify success
	setClient(processSuccessReponse)
	result, err := worker.ProcessMetadata(eviID, &payload)
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
		setClient(processErrorResponse_NoContent)
		result, err = worker.ProcessMetadata(eviID, &payload)
		verifyErrorScenario(result, err)
		require.NotNil(t, result.LastRunMessage)

		// with message
		setClient(processErrorResponse_WithMessage)
		result, err = worker.ProcessMetadata(eviID, &payload)
		verifyErrorScenario(result, err)
		require.Equal(t, content, *result.LastRunMessage)

		// without message
		setClient(processErrorResponse_StatusCode)
		result, err = worker.ProcessMetadata(eviID, &payload)
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
		setClient(processDeferalResponse_StatusCode)
		verifyDefferalResult(worker.ProcessMetadata(eviID, &payload))

		// action version
		setClient(processDeferalReponse_Action)
		verifyDefferalResult(worker.ProcessMetadata(eviID, &payload))
	}

	// verify Rejections
	{
		verifyRejectionResult := func(result *models.EvidenceMetadata, err error) {
			require.NoError(t, err)
			require.Equal(t, eviID, result.EvidenceID)
			require.Equal(t, evidencemetadata.StatusCompleted.Ptr(), result.Status)
		}
		// status code version
		setClient(processRejectedResponse_StatusCode)
		verifyRejectionResult(worker.ProcessMetadata(eviID, &payload))

		// action version
		setClient(processRejectedReponse_Action)
		result, err := worker.ProcessMetadata(eviID, &payload)
		verifyRejectionResult(result, err)
		require.Equal(t, content, *result.LastRunMessage)
	}
}
