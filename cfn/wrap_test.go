// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package cfn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

var testEvent = &Event{
	RequestType:       RequestCreate,
	RequestID:         "unique id for this create request",
	ResponseURL:       "http://pre-signed-S3-url-for-response",
	LogicalResourceID: "MyTestResource",
	StackID:           "arn:aws:cloudformation:us-west-2:EXAMPLE/stack-name/guid",
}

func TestCopyLambdaLogStream(t *testing.T) {
	lgs := lambdacontext.LogStreamName
	lambdacontext.LogStreamName = "DUMMYLOGSTREAMNAME"

	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			response := extractResponseBody(t, req)

			assert.Equal(t, StatusSuccess, response.Status)
			assert.Equal(t, testEvent.LogicalResourceID, response.LogicalResourceID)
			assert.Equal(t, "DUMMYLOGSTREAMNAME", response.PhysicalResourceID)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       nopCloser{bytes.NewBufferString("")},
			}, nil
		},
	}

	fn := func(ctx context.Context, event Event) (physicalResourceID string, data map[string]interface{}, err error) {
		return
	}

	lambdaWrapWithClient(fn, client)(context.TODO(), *testEvent)
	lambdacontext.LogStreamName = lgs
}

func TestDontCopyLogicalResourceId(t *testing.T) {
	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			response := extractResponseBody(t, req)

			assert.Equal(t, StatusSuccess, response.Status)
			assert.Equal(t, testEvent.LogicalResourceID, response.LogicalResourceID)
			assert.Equal(t, "testingtesting", response.PhysicalResourceID)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       nopCloser{bytes.NewBufferString("")},
			}, nil
		},
	}

	fn := func(ctx context.Context, event Event) (physicalResourceID string, data map[string]interface{}, err error) {
		physicalResourceID = "testingtesting"
		return
	}

	lambdaWrapWithClient(fn, client)(context.TODO(), *testEvent)
}

func TestWrappedError(t *testing.T) {
	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			response := extractResponseBody(t, req)

			assert.Equal(t, StatusFailed, response.Status)
			assert.Empty(t, response.PhysicalResourceID)
			assert.Equal(t, "failed to create resource", response.Reason)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       nopCloser{bytes.NewBufferString("")},
			}, nil
		},
	}

	fn := func(ctx context.Context, event Event) (physicalResourceID string, data map[string]interface{}, err error) {
		err = errors.New("failed to create resource")
		return
	}

	lambdaWrapWithClient(fn, client)(context.TODO(), *testEvent)
}

func TestWrappedSendFailure(t *testing.T) {
	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusForbidden,
			}, errors.New("things went wrong")
		},
	}

	fn := func(ctx context.Context, event Event) (physicalResourceID string, data map[string]interface{}, err error) {
		return
	}

	r, e := lambdaWrapWithClient(fn, client)(context.TODO(), *testEvent)
	assert.NotNil(t, e)
	assert.Equal(t, "things went wrong", r)
}

func extractResponseBody(t *testing.T, req *http.Request) Response {
	assert.NotContains(t, req.Header, "Content-Type")

	body, err := ioutil.ReadAll(req.Body)
	assert.NoError(t, err)
	var response Response
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	return response
}
