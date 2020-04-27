// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Function struct which wrap the Handler
type Function struct {
	handler Handler
	context context.Context
}

// NewFunction which creates a Function with a given Handler
func NewFunction(handler Handler) *Function {
	return &Function{handler: handler}
}

// NewFunctionWithContext which creates a Function with a given Handler and sets the base Context.
func NewFunctionWithContext(ctx context.Context, handler Handler) *Function {
	return &Function{
		context: ctx,
		handler: handler,
	}
}

// Context returns the base context used for the fn.
func (fn *Function) Context() context.Context {
	if fn.context == nil {
		return context.Background()
	}

	return fn.context
}

// Ping method which given a PingRequest and a PingResponse parses the PingResponse
func (fn *Function) Ping(req *messages.PingRequest, response *messages.PingResponse) error {
	*response = messages.PingResponse{}
	return nil
}

// Invoke method try to perform a command given an InvokeRequest and an InvokeResponse
func (fn *Function) Invoke(req *messages.InvokeRequest, response *messages.InvokeResponse) error {
	defer func() {
		if err := recover(); err != nil {
			panicInfo := getPanicInfo(err)
			response.Error = &messages.InvokeResponse_Error{
				Message:    panicInfo.Message,
				Type:       getErrorType(err),
				StackTrace: panicInfo.StackTrace,
				ShouldExit: true,
			}
		}
	}()

	deadline := time.Unix(req.Deadline.Seconds, req.Deadline.Nanos).UTC()
	invokeContext, cancel := context.WithDeadline(fn.Context(), deadline)
	defer cancel()

	lc := &lambdacontext.LambdaContext{
		AwsRequestID:       req.RequestId,
		InvokedFunctionArn: req.InvokedFunctionArn,
		Identity: lambdacontext.CognitoIdentity{
			CognitoIdentityID:     req.CognitoIdentityId,
			CognitoIdentityPoolID: req.CognitoIdentityPoolId,
		},
	}
	if len(req.ClientContext) > 0 {
		if err := json.Unmarshal(req.ClientContext, &lc.ClientContext); err != nil {
			response.Error = lambdaErrorResponse(err)
			return nil
		}
	}
	invokeContext = lambdacontext.NewContext(invokeContext, lc)

	invokeContext = context.WithValue(invokeContext, "x-amzn-trace-id", req.XAmznTraceId)
	os.Setenv("_X_AMZN_TRACE_ID", req.XAmznTraceId)

	payload, err := fn.handler.Invoke(invokeContext, req.Payload)
	if err != nil {
		response.Error = lambdaErrorResponse(err)
		return nil
	}
	response.Payload = payload
	return nil
}

func getErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}

func lambdaErrorResponse(invokeError error) *messages.InvokeResponse_Error {
	var errorName string
	if errorType := reflect.TypeOf(invokeError); errorType.Kind() == reflect.Ptr {
		errorName = errorType.Elem().Name()
	} else {
		errorName = errorType.Name()
	}
	return &messages.InvokeResponse_Error{
		Message: invokeError.Error(),
		Type:    errorName,
	}
}
