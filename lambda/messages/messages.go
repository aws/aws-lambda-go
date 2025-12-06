// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package messages

import "fmt"

type PingRequest struct {
}

type PingResponse struct {
}

//nolint:staticcheck
type InvokeRequest_Timestamp struct {
	Seconds int64
	Nanos   int64
}

//nolint:staticcheck
type InvokeRequest struct {
	Payload               []byte
	RequestId             string //nolint:staticcheck
	XAmznTraceId          string
	Deadline              InvokeRequest_Timestamp
	InvokedFunctionArn    string
	CognitoIdentityId     string //nolint:staticcheck
	CognitoIdentityPoolId string //nolint:staticcheck
	ClientContext         []byte
}

type InvokeResponse struct {
	Payload []byte
	Error   *InvokeResponse_Error
}

//nolint:staticcheck
type InvokeResponse_Error struct {
	Message    string                             `json:"errorMessage"`
	Type       string                             `json:"errorType"`
	StackTrace []*InvokeResponse_Error_StackFrame `json:"stackTrace,omitempty"`
	ShouldExit bool                               `json:"-"`
}

func (e InvokeResponse_Error) Error() string {
	return fmt.Sprintf("%#v", e)
}

//nolint:staticcheck
type InvokeResponse_Error_StackFrame struct {
	Path  string `json:"path"`
	Line  int32  `json:"line"`
	Label string `json:"label"`
}
