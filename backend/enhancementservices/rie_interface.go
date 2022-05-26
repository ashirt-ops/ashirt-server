package enhancementservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/theparanoids/ashirt-server/backend/helpers"
)

type LambdaRIEClient struct{}

func newRIELambdaClient() LambdaInvokableClient {

	return LambdaRIEClient{}
}

func (l LambdaRIEClient) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	if input.FunctionName == nil {
		return nil, fmt.Errorf("missing a function name for RIE lambda client")
	}
	url := "http://" + *input.FunctionName + ":8080/2015-03-31/functions/function/invocations"

	resp, err := helpers.MakeJSONRequest("POST", url, bytes.NewReader(input.Payload), helpers.NoMod)

	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	err = json.Unmarshal(respBody, &data)
	if err != nil {
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

	// this is all the aws_worker uses at the moment
	out := lambda.InvokeOutput{
		FunctionError: nil,
		Payload:       []byte(strBody),
		StatusCode:    helpers.Ptr(int64(resp.StatusCode)),
	}
	
	return &out, nil
}
