// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package cfn

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataCopiedFromRequest(t *testing.T) {
	e := &Event{
		RequestID:         "unique id for this create request",
		ResponseURL:       "http://pre-signed-S3-url-for-response",
		LogicalResourceID: "MyTestResource",
		StackID:           "arn:aws:cloudformation:us-west-2:EXAMPLE/stack-name/guid",
	}

	r := NewResponse(e)
	assert.Equal(t, e.RequestID, r.RequestID)
	assert.Equal(t, e.LogicalResourceID, r.LogicalResourceID)
	assert.Equal(t, e.StackID, r.StackID)
	assert.Equal(t, e.ResponseURL, r.url)
}

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{}, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestRequestSentCorrectly(t *testing.T) {
	r := &Response{
		Status: StatusSuccess,
		url:    "http://pre-signed-S3-url-for-response",
	}

	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.NotContains(t, req.Header, "Content-Type")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       nopCloser{bytes.NewBufferString("")},
			}, nil
		},
	}

	assert.NoError(t, r.sendWith(client))
}

func TestRequestForbidden(t *testing.T) {
	r := &Response{
		Status: StatusSuccess,
		url:    "http://pre-signed-S3-url-for-response",
	}

	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.NotContains(t, req.Header, "Content-Type")
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       nopCloser{bytes.NewBufferString("")},
			}, nil
		},
	}

	assert.Error(t, r.sendWith(client))
}
