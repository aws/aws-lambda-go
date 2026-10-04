// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package cfn

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"testing/iotest"

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

type trackingReadCloser struct {
	io.Reader
	closeCount int
	closeErr   error
}

func (r *trackingReadCloser) Close() error {
	r.closeCount++
	return r.closeErr
}

func TestResponseBodyClosed(t *testing.T) {
	readErr := errors.New("response body read failed")
	closeErr := errors.New("response body close failed")
	for _, test := range []struct {
		name       string
		statusCode int
		reader     io.Reader
		closeErr   error
		wantErr    error
	}{
		{"success", http.StatusOK, bytes.NewBufferString(""), nil, nil},
		{"HTTP error", http.StatusForbidden, bytes.NewBufferString("forbidden"), nil, fmt.Errorf("invalid status code. got: %d", http.StatusForbidden)},
		{"read error", http.StatusOK, iotest.ErrReader(readErr), nil, readErr},
		{"partial read error", http.StatusOK, io.MultiReader(bytes.NewBufferString("partial body"), iotest.ErrReader(readErr)), nil, readErr},
		{"close error", http.StatusOK, bytes.NewBufferString(""), closeErr, nil},
		{"read and close errors", http.StatusOK, iotest.ErrReader(readErr), closeErr, readErr},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := &trackingReadCloser{Reader: test.reader, closeErr: test.closeErr}
			client := &mockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: test.statusCode, Body: body}, nil
				},
			}
			r := &Response{Status: StatusSuccess, url: "http://pre-signed-S3-url-for-response"}

			assert.Equal(t, test.wantErr, r.sendWith(client))
			assert.Equal(t, 1, body.closeCount)
		})
	}
}

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

	sc := http.StatusForbidden
	client := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.NotContains(t, req.Header, "Content-Type")
			return &http.Response{
				StatusCode: sc,
				Body:       nopCloser{bytes.NewBufferString("")},
			}, nil
		},
	}

	s := r.sendWith(client)
	if assert.Error(t, s) {
		assert.Equal(t, fmt.Errorf("invalid status code. got: %d", sc), s)
	}
}
