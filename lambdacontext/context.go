// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

// Package lambdacontext provides access to Lambda execution context information.
//
// This package allows Lambda functions to access metadata about the current invocation,
// including request ID, function ARN, Cognito identity, and client context. Context
// information is retrieved from the standard Go context.Context using FromContext().
//
// See https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html
package lambdacontext

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strconv"
)

// W3CAllowedFields is the allowlist of W3C trace-context fields that may be
// surfaced through LambdaContext.W3C(). Any other key carried on
// clientContext.w3c is ignored, and any allowlisted key whose value is not a
// JSON string is dropped.
var W3CAllowedFields = []string{"traceparent", "tracestate", "baggage"}

// LogGroupName is the name of the log group that contains the log streams of the current Lambda Function
var LogGroupName string

// LogStreamName name of the log stream that the current Lambda Function's logs will be sent to
var LogStreamName string

// FunctionName the name of the current Lambda Function
var FunctionName string

// MemoryLimitInMB is the configured memory limit for the current instance of the Lambda Function
var MemoryLimitInMB int

// FunctionVersion is the published version of the current instance of the Lambda Function
var FunctionVersion string

var maxConcurrency int

func init() {
	LogGroupName = os.Getenv("AWS_LAMBDA_LOG_GROUP_NAME")
	LogStreamName = os.Getenv("AWS_LAMBDA_LOG_STREAM_NAME")
	FunctionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if limit, err := strconv.Atoi(os.Getenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE")); err != nil {
		MemoryLimitInMB = 0
	} else {
		MemoryLimitInMB = limit
	}
	FunctionVersion = os.Getenv("AWS_LAMBDA_FUNCTION_VERSION")
	if v, err := strconv.Atoi(os.Getenv("AWS_LAMBDA_MAX_CONCURRENCY")); err != nil || v < 1 {
		maxConcurrency = 1
	} else {
		maxConcurrency = v
	}
}

func MaxConcurrency() int {
	return maxConcurrency
}

// ClientApplication is metadata about the calling application.
type ClientApplication struct {
	InstallationID string `json:"installation_id"`
	AppTitle       string `json:"app_title"`
	AppVersionCode string `json:"app_version_code"`
	AppPackageName string `json:"app_package_name"`
}

// ClientContext is information about the client application passed by the calling application.
type ClientContext struct {
	Client ClientApplication
	Env    map[string]string `json:"env"`
	Custom map[string]string `json:"custom"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ClientContext.
// This handles the case where values in the "custom" map are not strings
// (e.g. nested JSON objects), by serializing non-string values back to
// their JSON string representation.
func (cc *ClientContext) UnmarshalJSON(data []byte) error {
	var raw struct {
		Client ClientApplication          `json:"Client"`
		Env    map[string]string          `json:"env"`
		Custom map[string]json.RawMessage `json:"custom"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	cc.Client = raw.Client
	cc.Env = raw.Env
	if raw.Custom != nil {
		cc.Custom = make(map[string]string, len(raw.Custom))
		for k, v := range raw.Custom {
			var s string
			if err := json.Unmarshal(v, &s); err == nil {
				cc.Custom[k] = s
			} else {
				cc.Custom[k] = string(v)
			}
		}
	}
	return nil
}

// CognitoIdentity is the cognito identity used by the calling application.
type CognitoIdentity struct {
	CognitoIdentityID     string
	CognitoIdentityPoolID string
}

// LambdaContext is the set of metadata that is passed for every Invoke.
type LambdaContext struct {
	AwsRequestID       string //nolint: staticcheck
	InvokedFunctionArn string //nolint: staticcheck
	Identity           CognitoIdentity
	ClientContext      ClientContext
	TenantID           string `json:",omitempty"`
	w3c                map[string]string
}

// W3C returns the W3C trace-context fields (see W3CAllowedFields) that were
// carried on clientContext.w3c at invoke time. The returned map is a fresh copy
// on every call, so mutating it never affects the context. It is never nil; an
// invoke that carried no W3C trace-context yields an empty map.
//
// The w3c key is not part of ClientContext, so it is deliberately never
// surfaced through lc.ClientContext — W3C() is the only accessor.
func (lc *LambdaContext) W3C() map[string]string {
	out := make(map[string]string, len(lc.w3c))
	for k, v := range lc.w3c {
		out[k] = v
	}
	return out
}

func (lc *LambdaContext) ExtractW3C(clientContextJSON []byte) {
	lc.w3c = extractW3CFields(clientContextJSON)
}

func extractW3CFields(clientContextJSON []byte) map[string]string {
	fields := map[string]string{}
	if len(clientContextJSON) == 0 {
		return fields
	}

	var envelope struct {
		W3C json.RawMessage `json:"w3c"`
	}
	if err := json.Unmarshal(clientContextJSON, &envelope); err != nil || len(envelope.W3C) == 0 {
		return fields
	}

	// w3c must be a JSON object; a string, array, number, etc. yields empty.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(envelope.W3C, &raw); err != nil {
		return fields
	}

	for _, key := range W3CAllowedFields {
		value, ok := raw[key]
		if !ok {
			continue
		}
		// Only keep values that are genuine JSON strings. A JSON string always
		// begins with a double-quote, so this rejects numbers, null, objects
		// and arrays without a second unmarshal attempt.
		trimmed := bytes.TrimSpace(value)
		if len(trimmed) == 0 || trimmed[0] != '"' {
			continue
		}
		var s string
		if err := json.Unmarshal(trimmed, &s); err == nil {
			fields[key] = s
		}
	}
	return fields
}

// An unexported type to be used as the key for types in this package.
// This prevents collisions with keys defined in other packages.
type key struct{}

// The key for a LambdaContext in Contexts.
// Users of this package must use lambdacontext.NewContext and lambdacontext.FromContext
// instead of using this key directly.
var contextKey = &key{}

// NewContext returns a new Context that carries value lc.
func NewContext(parent context.Context, lc *LambdaContext) context.Context {
	return context.WithValue(parent, contextKey, lc)
}

// FromContext returns the LambdaContext value stored in ctx, if any.
func FromContext(ctx context.Context) (*LambdaContext, bool) {
	lc, ok := ctx.Value(contextKey).(*LambdaContext)
	return lc, ok
}
