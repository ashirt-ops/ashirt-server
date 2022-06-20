// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices

import (
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
)

var allWorkers []string = []string{}
func AllWorkers() []string {
	return allWorkers
}

type BasicServiceWorkerConfig struct {
	Type    string `json:"type"`
	Version int64  `json:"version"`
}

type ServiceWorker interface {
	Build(workerName string, config []byte) error
	Test() ServiceTestResult
	Process(evidenceID int64, payload *NewEvidencePayload) (*models.EvidenceMetadata, error)
}

// ServiceTestResult provides a view of a Worker test
type ServiceTestResult struct {
	// Message contains helpful text detailing _why_ there was a failure
	Message string
	// Live indicates if the service is available or not
	Live bool
	// Error indicates if there was some fundamental error that prevented a full test
	Error error
}

type TestResp struct {
	Status  string  `json:"status"`
	Message *string `json:"message"`
}

func ErrorTestResult(err error) ServiceTestResult {
	return ServiceTestResult{
		Error: err,
		Live:  false,
	}
}

func ErrorTestResultWithMessage(err error, message string) ServiceTestResult {
	rtn := ErrorTestResult(err)
	rtn.Message = message
	return rtn
}

func TestResultSuccess(message string) ServiceTestResult {
	return ServiceTestResult{
		Message: message,
		Live:    true,
	}
}

type LambdaInvokableClient interface {
	Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

type ProcessResponse struct {
	Action  string  `json:"action"`  // Rejected | Deferred | Processed | Error
	Content *string `json:"content"` // Error => reason, Processed => Result
}

type RequestFn = func(method, url string, body io.Reader, updateRequest helpers.ModifyReqFunc) (*http.Response, error)
