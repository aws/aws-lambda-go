// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

//go:build !lambda.norpc
// +build !lambda.norpc

package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	srv := NewFunction(testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			if deadline, ok := ctx.Deadline(); ok {
				return deadline.UnixNano(), nil
			}
			return nil, errors.New("!?!?!?!?!")
		},
	))
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

func TestInvokeWithContext(t *testing.T) {
	// dummyKey creates a safe context to appease the linter
	type dummyKey struct{}
	var key dummyKey
	srv := NewFunction(&handlerOptions{
		handlerFunc: func(ctx context.Context, _ []byte) (io.Reader, error) {
			assert.Equal(t, "dummy", ctx.Value(key))
			if deadline, ok := ctx.Deadline(); ok {
				return strings.NewReader(strconv.FormatInt(deadline.UnixNano(), 10)), nil
			}
			return nil, errors.New("!?!?!?!?!")
		},
		baseContext: context.WithValue(context.Background(), key, "dummy"),
	})
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

func (e CustomError) Error() string { return "Something bad happened!" }

func TestCustomError(t *testing.T) {
	srv := NewFunction(testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			return nil, CustomError{}
		},
	))
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{}, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Payload)
	assert.Equal(t, "Something bad happened!", response.Error.Message)
	assert.Equal(t, "CustomError", response.Error.Type)
}

type CustomError2 struct{}

func (e *CustomError2) Error() string { return "Something bad happened!" }

func TestCustomErrorRef(t *testing.T) {

	srv := NewFunction(testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			return nil, &CustomError2{}
		},
	))
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{}, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Payload)
	assert.Equal(t, "Something bad happened!", response.Error.Message)
	assert.Equal(t, "CustomError2", response.Error.Type)
}

func TestContextPlumbing(t *testing.T) {
	srv := NewFunction(testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			lc, _ := lambdacontext.FromContext(ctx)
			return lc, nil
		},
	))
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

func TestXAmznTraceID(t *testing.T) {
	type XRayResponse struct {
		Env string
		Ctx string
	}

	srv := NewFunction(testWrapperHandler(
		func(ctx context.Context, input []byte) (interface{}, error) {
			return &XRayResponse{
				Env: os.Getenv("_X_AMZN_TRACE_ID"),
				Ctx: ctx.Value("x-amzn-trace-id").(string),
			}, nil
		},
	))

	sequence := []struct {
		Input    string
		Expected string
	}{
		{
			"",
			`{"Env": "", "Ctx": ""}`,
		},
		{
			"dummyid",
			`{"Env": "dummyid", "Ctx": "dummyid"}`,
		},
		{
			"",
			`{"Env": "", "Ctx": ""}`,
		},
		{
			"123dummyid",
			`{"Env": "123dummyid", "Ctx": "123dummyid"}`,
		},
		{
			"",
			`{"Env": "", "Ctx": ""}`,
		},
		{
			"",
			`{"Env": "", "Ctx": ""}`,
		},
		{
			"567",
			`{"Env": "567", "Ctx": "567"}`,
		},
		{
			"hihihi",
			`{"Env": "hihihi", "Ctx": "hihihi"}`,
		},
	}

	for i, test := range sequence {
		var response messages.InvokeResponse
		err := srv.Invoke(&messages.InvokeRequest{XAmznTraceId: test.Input}, &response)
		require.NoError(t, err, "failed test sequence[%d]", i)
		assert.JSONEq(t, test.Expected, string(response.Payload), "failed test sequence[%d]", i)
	}

}

type closeableResponse struct {
	reader io.Reader
	closed bool
}

func (c *closeableResponse) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *closeableResponse) Close() error {
	c.closed = true
	return nil
}

type readerError struct {
	err error
}

func (r *readerError) Read(_ []byte) (int, error) {
	return 0, r.err
}

func TestRPCModeInvokeClosesCloserIfResponseIsCloser(t *testing.T) {
	handlerResource := &closeableResponse{
		reader: strings.NewReader("<yolo/>"),
		closed: false,
	}
	srv := NewFunction(newHandler(func() (interface{}, error) {
		return handlerResource, nil
	}))
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{}, &response)
	require.NoError(t, err)
	assert.Equal(t, "<yolo/>", string(response.Payload))
	assert.True(t, handlerResource.closed)
}

func TestRPCModeInvokeReaderErrorPropogated(t *testing.T) {
	handlerResource := &closeableResponse{
		reader: &readerError{errors.New("yolo")},
		closed: false,
	}
	srv := NewFunction(newHandler(func() (interface{}, error) {
		return handlerResource, nil
	}))
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{}, &response)
	require.NoError(t, err)
	assert.Equal(t, "", string(response.Payload))
	assert.Equal(t, "yolo", response.Error.Message)
	assert.True(t, handlerResource.closed)
}
