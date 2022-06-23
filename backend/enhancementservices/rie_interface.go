// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/theparanoids/ashirt-server/backend/helpers"
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

	var data map[string]any
	if err = json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}
	body, ok := data["body"]
	if !ok {
		return nil, fmt.Errorf("unable to read lambda body")
	}
	strBody, ok := body.(string)
	if !ok {
		return nil, fmt.Errorf("lambda response body is not a string")
	}

	out.Payload = []byte(strBody)

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
