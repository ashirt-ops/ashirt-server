package enhancementservices

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaRIEClient struct {
	// makeRequestFn provides an alternative function to make a JSON based request. Should typically be nil,
	// except when unit testing
	makeRequestFn RequestFn
}

// lambdaMutex provides a lock on RIE calls. The RIE does not support parallelized calls.
var lambdaMutex sync.Mutex

func newRIELambdaClient() LambdaInvokableClient {
	return LambdaRIEClient{}
}

func MkRIEURL(lambdaName string) string {
	return "http://" + lambdaName + ":8080/2015-03-31/functions/function/invocations"
}

// Invoke mimics the aws Lambda function of the same name. This is useful for development testing
// without incurring AWS fees
func (l LambdaRIEClient) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	if input.FunctionName == nil {
		return nil, fmt.Errorf("missing a function name for RIE lambda client")
	}
	url := MkRIEURL(*input.FunctionName)

	lambdaMutex.Lock()
	resp, err := l.makeJSONRequest("POST", url, bytes.NewReader(input.Payload), helpers.NoMod)
	lambdaMutex.Unlock()

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := lambda.InvokeOutput{
		FunctionError: nil,
		StatusCode:    helpers.Ptr(int64(resp.StatusCode)),
	}
	if len(respBody) == 0 {
		return &out, nil
	}
	out.Payload = respBody

	return &out, nil
}

// NewTestRIELambdaClient creates an instance of LambdaInvokableClient that can be used for unit testing
func NewTestRIELambdaClient(fn RequestFn) LambdaInvokableClient {
	return LambdaRIEClient{
		makeRequestFn: fn,
	}
}

// makeJSONRequest is an abstraction over MakeJSONRequest to enable unit testing
func (l LambdaRIEClient) makeJSONRequest(method, url string, body io.Reader, updateRequest helpers.ModifyReqFunc) (*http.Response, error) {
	if l.makeRequestFn != nil {
		return l.makeRequestFn(method, url, body, updateRequest)
	}
	return helpers.MakeJSONRequest(method, url, body, updateRequest)
}
