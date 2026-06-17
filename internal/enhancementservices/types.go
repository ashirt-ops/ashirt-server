package enhancementservices

import (
	"io"
	"net/http"

	"github.com/ashirt-ops/ashirt-server/internal/helpers"
)

var allWorkers []string = []string{}

// AllWorkers is an effective constant representing all possible workers (not deleted)
func AllWorkers() []string {
	return allWorkers
}

type BasicServiceWorkerConfig struct {
	Type    string `json:"type"`
	Version int64  `json:"version"`
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

func errorTestResult(err error) ServiceTestResult {
	return ServiceTestResult{
		Error: err,
		Live:  false,
	}
}

func errorTestResultWithMessage(err error, message string) ServiceTestResult {
	rtn := errorTestResult(err)
	rtn.Message = message
	return rtn
}

func testResultSuccess(message string) ServiceTestResult {
	return ServiceTestResult{
		Message: message,
		Live:    true,
	}
}

type ProcessResponse struct {
	Action  string  `json:"action"`  // Rejected | Deferred | Processed | Error
	Content *string `json:"content"` // Error => reason, Processed => Result
}

type RequestFn = func(method, url string, body io.Reader, updateRequest helpers.ModifyReqFunc) (*http.Response, error)
