// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Function struct which wrap the Handler
//
// Deprecated: The Function type is public for the go1.x runtime internal use of the net/rpc package
type Function struct {
	handler *handlerOptions
}

// NewFunction which creates a Function with a given Handler
//
// Deprecated: The Function type is public for the go1.x runtime internal use of the net/rpc package
func NewFunction(handler Handler) *Function {
	return &Function{newHandler(handler)}
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
			response.Error = lambdaPanicResponse(err)
		}
	}()

	deadline := time.Unix(req.Deadline.Seconds, req.Deadline.Nanos).UTC()
	invokeContext, cancel := context.WithDeadline(fn.baseContext(), deadline)
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

	// nolint:staticcheck
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

func (fn *Function) baseContext() context.Context {
	if fn.handler.baseContext != nil {
		return fn.handler.baseContext
	}
	return context.Background()
}
