package events

import (
	"encoding/json"
	"io/ioutil" //nolint: staticcheck
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestAppSyncIdentity_IAM(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("./testdata/appsync-identity-iam.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputIdentity AppSyncIAMIdentity
	if err = json.Unmarshal(inputJSON, &inputIdentity); err != nil {
		t.Errorf("could not unmarshal identity. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputIdentity)
	if err != nil {
		t.Errorf("could not marshal identity. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestAppSyncIdentity_Cognito(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("./testdata/appsync-identity-cognito.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputIdentity AppSyncCognitoIdentity
	if err = json.Unmarshal(inputJSON, &inputIdentity); err != nil {
		t.Errorf("could not unmarshal identity. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputIdentity)
	if err != nil {
		t.Errorf("could not marshal identity. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestAppSyncLambdaAuthorizerRequestMarshalling(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("./testdata/appsync-lambda-auth-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputEvent AppSyncLambdaAuthorizerRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestAppSyncLambdaAuthorizerRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, AppSyncLambdaAuthorizerRequest{})
}

func TestAppSyncLambdaAuthorizerResponseMarshalling(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("./testdata/appsync-lambda-auth-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputEvent AppSyncLambdaAuthorizerResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestAppSyncLambdaAuthorizerResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, AppSyncLambdaAuthorizerResponse{})
}
