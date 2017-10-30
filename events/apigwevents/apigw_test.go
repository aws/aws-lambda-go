// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package apigwevents

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestApiGatewayRequestMarshaling(t *testing.T) {

	// read json from file
	inputJson, err := ioutil.ReadFile("./test-files/apigw-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent ApiGatewayProxyRequest
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJson, outputJson)
}

func TestApiGatewayRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ApiGatewayProxyRequest{})
}

func TestApiGatewayResponseMarshaling(t *testing.T) {

	// read json from file
	inputJson, err := ioutil.ReadFile("./test-files/apigw-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent ApiGatewayProxyResponse
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJson, outputJson)
}

func TestApiGatewayResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ApiGatewayProxyResponse{})
}
