// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package cfn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil" //nolint: staticcheck
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testEvent = &Event{
	RequestType:        RequestUpdate,
	RequestID:          "unique id for this create request",
	ResponseURL:        "http://pre-signed-S3-url-for-response",
	LogicalResourceID:  "MyTestResource",
	PhysicalResourceID: "prevPhysicalResourceID",
	StackID:            "arn:aws:cloudformation:us-west-2:EXAMPLE/stack-name/guid",
}

func TestLambdaPhysicalResourceId(t *testing.T) {

	tests := []struct {
		// Input to the lambda
		inputRequestType RequestType

		// Output from the lambda
		returnErr                error
		returnPhysicalResourceID string

		// The PhysicalResourceID to test
		expectedPhysicalResourceID string
	}{
		// For Create with no returned PhysicalResourceID
		{RequestCreate, nil, "", testEvent.RequestID}, // Use RequestID as default for PhysicalResourceID
		{RequestCreate, fmt.Errorf("dummy error"), "", testEvent.RequestID},

		// For Create with PhysicalResourceID
		{RequestCreate, nil, "newPhysicalResourceID", "newPhysicalResourceID"},
		{RequestCreate, fmt.Errorf("dummy error"), "newPhysicalResourceID", "newPhysicalResourceID"},

		// For Update with no returned PhysicalResourceID
		{RequestUpdate, nil, "", "prevPhysicalResourceID"},
		{RequestUpdate, fmt.Errorf("dummy error"), "", "prevPhysicalResourceID"},

		// For Update with returned PhysicalResourceID
		{RequestUpdate, nil, "newPhysicalResourceID", "newPhysicalResourceID"},
		{RequestUpdate, fmt.Errorf("dummy error"), "newPhysicalResourceID", "newPhysicalResourceID"},

		// For Delete with no returned PhysicalResourceID
		{RequestDelete, nil, "", "prevPhysicalResourceID"},
		{RequestDelete, fmt.Errorf("dummy error"), "", "prevPhysicalResourceID"},

		// For Delete with returned PhysicalResourceID = old PhysicalResourceID
		{RequestDelete, nil, "prevPhysicalResourceID", "prevPhysicalResourceID"},
		{RequestDelete, fmt.Errorf("dummy error"), "prevPhysicalResourceID", "prevPhysicalResourceID"},

		// For Delete with returned PhysicalResourceID != old PhysicalResourceID
		// Technically, lambda handlers shouldn't return a different physical resource id upon deletion.
		// Typescript CDK implementation for handlers would return an error here.
		// But CFn ignores the returned PhysicalResourceID for deletion so it doesn't matter.
		{RequestDelete, nil, "newPhysicalResourceID", "newPhysicalResourceID"},
		{RequestDelete, fmt.Errorf("dummy error"), "newPhysicalResourceID", "newPhysicalResourceID"},
	}
	for _, test := range tests {

		curTestEvent := *testEvent
		curTestEvent.RequestType = test.inputRequestType

		if curTestEvent.RequestType == RequestCreate {
			curTestEvent.PhysicalResourceID = ""
		}

		client := &mockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				response := extractResponseBody(t, req)

				if test.returnErr == nil {
					assert.Equal(t, StatusSuccess, response.Status)
				} else {
					assert.Equal(t, StatusFailed, response.Status)
				}

				assert.Equal(t, curTestEvent.LogicalResourceID, response.LogicalResourceID)
				assert.Equal(t, test.expectedPhysicalResourceID, response.PhysicalResourceID)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       nopCloser{bytes.NewBufferString("")},
				}, nil
			},
		}

		fn := func(ctx context.Context, event Event) (physicalResourceID string, data map[string]interface{}, err error) {
			return test.returnPhysicalResourceID, nil, test.returnErr
		}

		_, err := lambdaWrapWithClient(fn, client)(context.TODO(), curTestEvent)
		assert.NoError(t, err)
	}
}

func TestPanicSendsFailure(t *testing.T) {
	didSendStatus := false

	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			response := extractResponseBody(t, req)
			assert.Equal(t, StatusFailed, response.Status)

			// Even in a panic, a dummy PhysicalResourceID should be returned to ensure the error is surfaced correctly
			assert.Equal(t, testEvent.PhysicalResourceID, response.PhysicalResourceID)
			didSendStatus = response.Status == StatusFailed

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       nopCloser{bytes.NewBufferString("")},
			}, nil
		},
	}

	fn := func(ctx context.Context, event Event) (physicalResourceID string, data map[string]interface{}, err error) {
		err = errors.New("some panic that shouldn't be caught")
		panic(err)
	}

	assert.Panics(t, func() {
		_, err := lambdaWrapWithClient(fn, client)(context.TODO(), *testEvent)
		assert.NoError(t, err)
	})

	assert.True(t, didSendStatus, "FAILED should be sent to CloudFormation service")
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

	_, err := lambdaWrapWithClient(fn, client)(context.TODO(), *testEvent)
	assert.NoError(t, err)
}

func TestWrappedError(t *testing.T) {
	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			response := extractResponseBody(t, req)

			assert.Equal(t, StatusFailed, response.Status)
			assert.Equal(t, testEvent.PhysicalResourceID, response.PhysicalResourceID)
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

	_, err := lambdaWrapWithClient(fn, client)(context.TODO(), *testEvent)
	assert.NoError(t, err)
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
