package enhancementservices

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/servicetypes/evidencemetadata"
)

var lambdaClient *lambda.Lambda = nil

type awsConfigV1Worker struct {
	Config     AwsConfigV1
	WorkerName string
}

type AwsConfigV1 struct {
	BasicServiceWorkerConfig
	LambdaName string `json:"lambdaName"`
	AsyncFn    bool   `json:"asyncFunction"`
}

type awsTestResp struct {
	Status  string  `json:"status"`
	Message *string `json:"message"`
}

type awsProcessResp struct {
	Action  string  `json:"action"`  // Rejected | Deferred | Processed | Error
	Content *string `json:"content"` // Error => reason, Processed => Result
}

func buildLambdaClient() error {
	sess, err := session.NewSession()
	if err != nil {
		return backend.WrapError("Unable to establish an aws lambda session", err)
	}
	creds := credentials.NewCredentials(&credentials.StaticProvider{
		Value: credentials.Value{
			// TODO
		},
	})
	lambdaClient = lambda.New(sess, &aws.Config{
		Region:      helpers.Ptr(config.AWSRegion()),
		Credentials: creds,
	})

	return nil
}

func (w *awsConfigV1Worker) Build(workerName string, workerConfig []byte) error {
	if lambdaClient == nil {
		if err := buildLambdaClient(); err != nil {
			return err
		}
	}

	var awsConfig AwsConfigV1
	if err := json.Unmarshal([]byte(workerConfig), &awsConfig); err != nil {
		return err // TODO
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
	if w.Config.AsyncFn {
		input.SetInvocationType("Event")
	}

	out, err := lambdaClient.Invoke(&input)

	if err != nil {
		return ErrorTestResultWithMessage(err, "Unable to verify worker status")
	}

	if out.FunctionError != nil {
		return ErrorTestResultWithMessage(err, "Service experienced an error: "+*out.FunctionError)
	}

	var parsedData awsTestResp
	if err := json.Unmarshal(out.Payload, &parsedData); err != nil {
		return ErrorTestResultWithMessage(err, "Unable to parse response")
	}

	if parsedData.Status == "ok" {
		return TestResultSuccess("Service is functional")
	}
	if parsedData.Status == "error" {
		if parsedData.Message != nil {
			return ErrorTestResultWithMessage(err, "Service reported an error: "+*parsedData.Message)
		} else {
			return ErrorTestResultWithMessage(err, "Service reported an error")
		}
	}

	return ErrorTestResultWithMessage(nil, "Service did not reply with a supported status")
}

func (w *awsConfigV1Worker) Process(evidenceID int64, payload *Payload) (*models.EvidenceMetadata, error) {
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

func handleAWSProcessResponse(dbModel *models.EvidenceMetadata, output *lambda.InvokeOutput) {
	recordRejection := func(message *string) {
		dbModel.Status = evidencemetadata.StatusCompleted.Ptr()
		dbModel.CanProcess = helpers.Ptr(false)
		dbModel.LastRunMessage = message
	}
	recordError := func(message *string) {
		dbModel.Status = evidencemetadata.StatusError.Ptr()
		dbModel.LastRunMessage = message
	}
	recordDeferral := func() {
		dbModel.Status = evidencemetadata.StatusQueued.Ptr()
		dbModel.CanProcess = helpers.Ptr(true)
	}
	recordProcessed := func(content string) {
		dbModel.Status = evidencemetadata.StatusCompleted.Ptr()
		dbModel.CanProcess = helpers.Ptr(true)
		dbModel.Body = content
	}

	statusCode := *output.StatusCode
	var parsedData awsProcessResp
	err := json.Unmarshal(output.Payload, &parsedData)

	if err != nil {
		recordError(helpers.Ptr("Unable to parse response"))
		return
	}

	switch statusCode {
	case http.StatusOK: // 200
		switch parsedData.Action {
		case "processed":
			if parsedData.Content != nil {
				recordProcessed(*parsedData.Content)
			} else {
				recordError(helpers.Ptr("Content was not delivered for successful run"))
			}
		case "rejected":
			recordRejection(parsedData.Content)
		case "error":
			recordError(parsedData.Content)
		case "deferred":
			recordDeferral()
		default:
			recordError(helpers.SprintfPtr("Unexpected response format (%v)", parsedData.Action))
		}
	case http.StatusAccepted:
		recordDeferral()
	case http.StatusNotAcceptable:
		recordRejection(nil)
	case http.StatusInternalServerError:
		recordError(nil)
	default:
		recordError(helpers.SprintfPtr("Unexpected response status code (%v)", statusCode))
	}
}
