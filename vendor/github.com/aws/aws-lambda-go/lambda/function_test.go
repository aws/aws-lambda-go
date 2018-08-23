// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

type testWrapperHandler func(ctx context.Context, input []byte) (interface{}, error)

func (h testWrapperHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	response, err := h(ctx, payload)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

// verify testWrapperHandler implements Handler
var _ Handler = (testWrapperHandler)(nil)

func TestInvoke(t *testing.T) {
	srv := &Function{handler: testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			if deadline, ok := ctx.Deadline(); ok {
				return deadline.UnixNano(), nil
			}
			return nil, errors.New("!?!?!?!?!")
		},
	)}
	deadline := time.Now()
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: deadline.Unix(),
			Nanos:   int64(deadline.Nanosecond()),
		}}, &response)
	assert.NoError(t, err)
	var responseValue int64
	assert.NoError(t, json.Unmarshal(response.Payload, &responseValue))
	assert.Equal(t, deadline.UnixNano(), responseValue)
}

type CustomError struct{}

func (e CustomError) Error() string { return fmt.Sprintf("Something bad happened!") }

func TestCustomError(t *testing.T) {

	srv := &Function{handler: testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			return nil, CustomError{}
		},
	)}
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{}, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Payload)
	assert.Equal(t, "Something bad happened!", response.Error.Message)
	assert.Equal(t, "CustomError", response.Error.Type)
}

type CustomError2 struct{}

func (e *CustomError2) Error() string { return fmt.Sprintf("Something bad happened!") }

func TestCustomErrorRef(t *testing.T) {

	srv := &Function{handler: testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			return nil, &CustomError2{}
		},
	)}
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{}, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Payload)
	assert.Equal(t, "Something bad happened!", response.Error.Message)
	assert.Equal(t, "CustomError2", response.Error.Type)
}

func TestContextPlumbing(t *testing.T) {
	srv := &Function{handler: testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			lc, _ := lambdacontext.FromContext(ctx)
			return lc, nil
		},
	)}
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{
		CognitoIdentityId:     "dummyident",
		CognitoIdentityPoolId: "dummypool",
		ClientContext: []byte(`{
			"Client": {
				"app_title": "dummytitle",
				"installation_id": "dummyinstallid",
				"app_version_code": "dummycode",
				"app_package_name": "dummyname"
			}
		}`),
		RequestId:          "dummyid",
		InvokedFunctionArn: "dummyarn",
	}, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Payload)
	expected := `
	{
		"AwsRequestID": "dummyid",
		"InvokedFunctionArn": "dummyarn",
		"Identity": {
			"CognitoIdentityID": "dummyident",
			"CognitoIdentityPoolID": "dummypool"
		},
		"ClientContext": {
			"Client": {
				"installation_id": "dummyinstallid",
				"app_title": "dummytitle",
				"app_version_code": "dummycode",
				"app_package_name": "dummyname"
			},
			"env": null,
			"custom": null
		}
	}
	`
	assert.JSONEq(t, expected, string(response.Payload))
}
