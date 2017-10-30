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
	"github.com/stretchr/testify/assert"
)

func TestInvoke(t *testing.T) {
	srv := &Function{handler: func(ctx context.Context, input []byte) (interface{}, error) {
		if deadline, ok := ctx.Deadline(); ok {
			return deadline.UnixNano(), nil
		}
		return nil, errors.New("!?!?!?!?!")
	}}
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

	srv := &Function{handler: func(ctx context.Context, input []byte) (interface{}, error) {
		return nil, CustomError{}
	}}
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

	srv := &Function{handler: func(ctx context.Context, input []byte) (interface{}, error) {
		return nil, &CustomError2{}
	}}
	var response messages.InvokeResponse
	err := srv.Invoke(&messages.InvokeRequest{}, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Payload)
	assert.Equal(t, "Something bad happened!", response.Error.Message)
	assert.Equal(t, "CustomError2", response.Error.Type)
}
