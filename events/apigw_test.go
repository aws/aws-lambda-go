// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
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

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
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

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
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

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
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

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestApiGatewayCustomAuthorizerRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayCustomAuthorizerRequest{})
}

func TestApiGatewayCustomAuthorizerRequestTypeRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayCustomAuthorizerRequestTypeRequest{})
}

func TestApiGatewayWebsocketRequestMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-websocket-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayWebsocketProxyRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestApiGatewayWebsocketRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayWebsocketProxyRequest{})
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

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestApiGatewayCustomAuthorizerResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, APIGatewayCustomAuthorizerResponse{})
}

func TestApiGatewayRestApiOpenApiRequestMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-restapi-openapi-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayProxyRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// validate request context
	requestContext := inputEvent.RequestContext
	if requestContext.OperationName != "HelloWorld" {
		t.Errorf("could not extract operationName from context: %v", requestContext)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestApiGatewayV2HTTPRequestJWTAuthorizerMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-v2-request-jwt-authorizer.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayV2HTTPRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// validate custom authorizer context
	authContext := inputEvent.RequestContext.Authorizer
	if authContext.JWT.Claims["claim1"] != "value1" {
		t.Errorf("could not extract authorizer claim from JWT: %v", authContext)
	}

	// validate HTTP details
	http := inputEvent.RequestContext.HTTP
	if http.Path != "/my/path" {
		t.Errorf("could not extract HTTP details: %v", http)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestApiGatewayV2HTTPRequestLambdaAuthorizerMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-v2-request-lambda-authorizer.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayV2HTTPRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// validate custom authorizer context
	authContext := inputEvent.RequestContext.Authorizer
	if authContext.Lambda["key"] != "value" {
		t.Errorf("could not extract authorizer information from Lambda authorizer: %v", authContext)
	}

	// validate HTTP details
	http := inputEvent.RequestContext.HTTP
	if http.Path != "/my/path" {
		t.Errorf("could not extract HTTP details: %v", http)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestApiGatewayV2HTTPRequestIAMMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-v2-request-iam.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayV2HTTPRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// validate custom authorizer context
	authContext := inputEvent.RequestContext.Authorizer
	if authContext.IAM.AccessKey != "ARIA2ZJZYVUEREEIHAKY" {
		t.Errorf("could not extract accessKey from IAM authorizer: %v", authContext)
	}
	if authContext.IAM.AccountID != "1234567890" {
		t.Errorf("could not extract accountId from IAM authorizer: %v", authContext)
	}
	if authContext.IAM.CallerID != "AROA7ZJZYVRE7C3DUXHH6:CognitoIdentityCredentials" {
		t.Errorf("could not extract callerId from IAM authorizer: %v", authContext)
	}
	if authContext.IAM.CognitoIdentity.AMR[0] != "foo" {
		t.Errorf("could not extract amr from CognitoIdentity: %v", authContext)
	}
	if authContext.IAM.CognitoIdentity.IdentityID != "us-east-1:3f291106-8703-466b-8f2b-3ecee1ca56ce" {
		t.Errorf("could not extract identityId from CognitoIdentity: %v", authContext)
	}
	if authContext.IAM.CognitoIdentity.IdentityPoolID != "us-east-1:4f291106-8703-466b-8f2b-3ecee1ca56ce" {
		t.Errorf("could not extract identityPoolId from CognitoIdentity: %v", authContext)
	}
	if authContext.IAM.PrincipalOrgID != "AwsOrgId" {
		t.Errorf("could not extract principalOrgId from IAM authorizer: %v", authContext)
	}
	if authContext.IAM.UserARN != "arn:aws:iam::1234567890:user/Admin" {
		t.Errorf("could not extract userArn from IAM authorizer: %v", authContext)
	}
	if authContext.IAM.UserID != "AROA2ZJZYVRE7Y3TUXHH6" {
		t.Errorf("could not extract userId from IAM authorizer: %v", authContext)
	}

	// validate HTTP details
	http := inputEvent.RequestContext.HTTP
	if http.Path != "/my/path" {
		t.Errorf("could not extract HTTP details: %v", http)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestApiGatewayV2HTTPRequestNoAuthorizerMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/apigw-v2-request-no-authorizer.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent APIGatewayV2HTTPRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// validate custom authorizer context
	authContext := inputEvent.RequestContext.Authorizer
	if authContext != nil {
		t.Errorf("unexpected authorizer: %v", authContext)
	}

	// validate HTTP details
	http := inputEvent.RequestContext.HTTP
	if http.Path != "/" {
		t.Errorf("could not extract HTTP details: %v", http)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}
