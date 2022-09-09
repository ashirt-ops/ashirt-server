// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
)

var lambdaClient LambdaInvokableClient = nil

type awsConfigV1Worker struct {
	Config     AWSConfigV1
	WorkerName string
}

type AWSConfigV1 struct {
	BasicServiceWorkerConfig
	LambdaName string `json:"lambdaName"`
	AsyncFn    bool   `json:"asyncFunction"`
}

func buildLambdaClient() error {
	sess, err := session.NewSession()
	if err != nil {
		return backend.WrapError("unable to establish an aws lambda session", err)
	}
	if config.UseLambdaRIE() {
		lambdaClient = newRIELambdaClient()
	} else {
		lambdaClient = lambda.New(sess, &aws.Config{
			Region: helpers.Ptr(config.AWSRegion()),
		})
	}

	return nil
}

func (w *awsConfigV1Worker) Build(workerName string, workerConfig []byte) error {
	// Create long-running lambda client, since we don't need a new one for each worker
	if lambdaClient == nil {
		if err := buildLambdaClient(); err != nil {
			return err
		}
	}

	var awsConfig AWSConfigV1
	if err := json.Unmarshal([]byte(workerConfig), &awsConfig); err != nil {
		return backend.WrapError("aws worker config is unparsable", err)
	}
	w.WorkerName = workerName
	w.Config = awsConfig
	return nil
}

func (w *awsConfigV1Worker) Test() ServiceTestResult {
	input := lambda.InvokeInput{
		FunctionName: &w.Config.LambdaName,
		Payload:      []byte(`{"type": "test"}`),
	}

	out, err := lambdaClient.Invoke(&input)

	if err != nil {
		return errorTestResultWithMessage(err, "Unable to verify worker status")
	}

	if out.FunctionError != nil {
		return errorTestResultWithMessage(nil, "Service experienced an error: "+*out.FunctionError)
	}

	var lambdaResponse LambdaResponse
	if err := json.Unmarshal(out.Payload, &lambdaResponse); err != nil {
		return errorTestResultWithMessage(err, "Unable to parse response")
	}

	if lambdaResponse.StatusCode != 200 {
		return errorTestResultWithMessage(err, "Lambda failed")
	}

	var parsedData TestResp
	if err := json.Unmarshal([]byte(lambdaResponse.Body), &parsedData); err != nil {
		return errorTestResultWithMessage(err, "Unable to parse response")
	}

	if parsedData.Status == "ok" {
		return testResultSuccess("Service is functional")
	}
	if parsedData.Status == "error" {
		if parsedData.Message != nil {
			return errorTestResultWithMessage(nil, "Service reported an error: "+*parsedData.Message)
		}
		return errorTestResultWithMessage(nil, "Service reported an error")
	}

	return errorTestResultWithMessage(nil, "Service did not reply with a supported status")
}

func (w *awsConfigV1Worker) ProcessMetadata(evidenceID int64, payload *NewEvidencePayload) (*models.EvidenceMetadata, error) {
	body, err := json.Marshal(*payload)
	if err != nil {
		return nil, backend.WrapError("unable to construct body", err)
	}

	input := lambda.InvokeInput{
		FunctionName: &w.Config.LambdaName,
		Payload:      body,
	}
	if w.Config.AsyncFn {
		input.SetInvocationType("Event")
	}

	out, err := lambdaClient.Invoke(&input)
	if err != nil {
		return nil, backend.WrapError("Unable to invoke lambda function", err)
	}
	if out.FunctionError != nil {
		return nil, backend.WrapError("Lambda invocation failed", err)
	}

	// handle deferral -- we can assume this is true if we set the invocation type to "event"
	model := models.EvidenceMetadata{
		Source:     w.WorkerName,
		EvidenceID: evidenceID,
	}

	handleAWSProcessResponse(&model, out)
	return &model, nil
}

func (w *awsConfigV1Worker) ProcessEvent(payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return backend.WrapError("unable to construct body", err)
	}

	input := lambda.InvokeInput{
		FunctionName: &w.Config.LambdaName,
		Payload:      body,
	}
	if w.Config.AsyncFn {
		input.SetInvocationType("Event")
	}

	if out, err := lambdaClient.Invoke(&input); err != nil {
		return backend.WrapError("Unable to invoke lambda function", err)
	} else if out.FunctionError != nil {
		return backend.WrapError("Lambda invocation failed", err)
	}

	return nil
}

func handleAWSProcessResponse(dbModel *models.EvidenceMetadata, output *lambda.InvokeOutput) {
	statusCode := *output.StatusCode
	var parsedData ProcessResponse

	if len(output.Payload) > 0 {
		var lambdaResponse LambdaResponse
		if err := json.Unmarshal(output.Payload, &lambdaResponse); err != nil {
			return
		}

		if err := json.Unmarshal([]byte(lambdaResponse.Body), &parsedData); err != nil {
			recordError(dbModel, helpers.Ptr("Unable to parse response"))
			return
		}
	}

	handleProcessResponse(dbModel, int(statusCode), parsedData)
}

// SetTestLambdaClient provides a way to conduct unit tests. Not intended for regular use
func SetTestLambdaClient(client LambdaInvokableClient) {
	lambdaClient = client
}

// BuildTestLambdaWorker provides a way to conduct unit tests.
// This function creates a canned worker suitable for immediate use.
// Not intended for regular use
func BuildTestLambdaWorker() awsConfigV1Worker {
	return BuildTestLambdaWorkerWithName("test-worker")
}

// BuildTestLambdaWorkerWithName provides a way to conduct unit tests.
// This function creates a canned worker suitable for immediate use, of the provided name.
// Not intended for regular use
func BuildTestLambdaWorkerWithName(name string) awsConfigV1Worker {
	return awsConfigV1Worker{
		WorkerName: name,
		Config: AWSConfigV1{
			AsyncFn:    false,
			LambdaName: name,
			BasicServiceWorkerConfig: BasicServiceWorkerConfig{
				Type:    "aws",
				Version: 1,
			},
		},
	}
}
