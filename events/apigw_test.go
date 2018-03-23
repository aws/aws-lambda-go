// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestApiGatewayRequestMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayProxyRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// validate custom authorizer context
	authContext := inputEvent.RequestContext.Authorizer
	if authContext["principalId"] != "admin" ||
		authContext["clientId"] != 1.0 ||
		authContext["clientName"] != "Exata" {
		t.Errorf("could not extract authorizer context: %v", authContext)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestApiGatewayRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayProxyRequest{})
}

func TestApiGatewayResponseMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayProxyResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestApiGatewayResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayProxyResponse{})
}

func TestApiGatewayCustomAuthorizerRequestMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-custom-auth-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayCustomAuthorizerRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestApiGatewayCustomAuthorizerRequestTypeRequestMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-custom-auth-request-type-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayCustomAuthorizerRequestTypeRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestApiGatewayCustomAuthorizerRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayCustomAuthorizerRequest{})
}

func TestApiGatewayCustomAuthorizerRequestTypeRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayCustomAuthorizerRequestTypeRequest{})
}

func TestApiGatewayCustomAuthorizerResponseMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-custom-auth-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayCustomAuthorizerResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestApiGatewayCustomAuthorizerResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayCustomAuthorizerResponse{})
}
