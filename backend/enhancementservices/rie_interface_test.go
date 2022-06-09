// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/enhancementservices"
	"github.com/theparanoids/ashirt-server/backend/helpers"
)

func TestInvoke(t *testing.T) {
	// variables to tweak the response/tests
	lambdaName := "magic"
	body := LambdaResponseData{}
	expectError := false

	mockInput := RequestMock{
		OnInvoked: func(rd RequestData) {
			actualUrl, err := url.Parse(rd.URL)
			require.NoError(t, err)
			require.Equal(t, lambdaName+":8080", actualUrl.Host)

			bodyBytes, _ := json.Marshal(body)
			require.Equal(t, bodyBytes, rd.Body)
			if expectError {
				require.Error(t, rd.Error)
			} else {
				require.NoError(t, rd.Error)
			}
		},
		OnSendRequest: func(req *http.Request, err error) (*http.Response, error) {
			require.NoError(t, err)
			w := httptest.NewRecorder()
			w.WriteHeader(http.StatusAccepted)
			bodyBytes, _ := json.Marshal(body)
			w.Write(bodyBytes)
			return w.Result(), nil
		},
	}

	client := enhancementservices.NewTestRIELambdaClient(makeMockRequestHandler(mockInput))

	// verify error
	_, err := client.Invoke(&lambda.InvokeInput{
		FunctionName: nil,
	})
	require.Error(t, err)

	expectedBody := `{"status":"ok"}`
	body.Body = expectedBody
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)
	out, err := client.Invoke(&lambda.InvokeInput{
		FunctionName: &lambdaName,
		Payload:      bodyBytes,
	})
	strPayload := string(out.Payload)
	fmt.Println(strPayload)
	require.NoError(t, err)
	require.Equal(t, []byte(expectedBody), out.Payload)
}

func makeMockRequestHandler(mock RequestMock) enhancementservices.RequestFn {
	return func(method, url string, body io.Reader, updateRequest helpers.ModifyReqFunc) (*http.Response, error) {
		content, err := io.ReadAll(body)
		clonedBody := bytes.NewReader(content)
		req := httptest.NewRequest(method, url, clonedBody)

		if mock.OnInvoked != nil {
			mock.OnInvoked(RequestData{
				Method:  method,
				URL:     url,
				Body:    content,
				Request: req,
				Error:   err,
			})
		}
		err = updateRequest(req)
		if mock.OnSendRequest != nil {
			return mock.OnSendRequest(req, err)
		}

		// default in case someone doesn't provide a RespondWith function
		w := httptest.NewRecorder()
		w.WriteHeader(http.StatusNoContent)
		return w.Result(), nil
	}
}

// opting for a struct here so On* functions can be omitted
type RequestMock struct {
	OnInvoked     func(RequestData)
	OnSendRequest func(*http.Request, error) (*http.Response, error)
}

type RequestData struct {
	Method  string
	URL     string
	Body    []byte
	Request *http.Request
	Error   error
}

type LambdaResponseData struct {
	Body string `json:"body"`
}
