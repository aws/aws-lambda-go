// Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"errors"
	"github.com/segmentio/encoding/json"
	"io/ioutil" //nolint: staticcheck
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLambdaFunctionURLResponseMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/lambda-urls-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent LambdaFunctionURLResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestLambdaFunctionURLRequestMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/lambda-urls-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into Go object
	var inputEvent LambdaFunctionURLRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestLambdaFunctionURLStreamingResponseMarshaling(t *testing.T) {
	for _, test := range []struct {
		name         string
		response     *LambdaFunctionURLStreamingResponse
		expectedHead string
		expectedBody string
	}{
		{
			"empty",
			&LambdaFunctionURLStreamingResponse{},
			`{"statusCode":200}`,
			"",
		},
		{
			"just the status code",
			&LambdaFunctionURLStreamingResponse{
				StatusCode: http.StatusTeapot,
			},
			`{"statusCode":418}`,
			"",
		},
		{
			"status and headers and cookies and body",
			&LambdaFunctionURLStreamingResponse{
				StatusCode: http.StatusTeapot,
				Headers:    map[string]string{"hello": "world"},
				Cookies:    []string{"cookies", "are", "yummy"},
				Body:       strings.NewReader(`<html>Hello Hello</html>`),
			},
			`{"statusCode":418, "headers":{"hello":"world"}, "cookies":["cookies","are","yummy"]}`,
			`<html>Hello Hello</html>`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			response, err := ioutil.ReadAll(test.response)
			require.NoError(t, err)
			sep := "\x00\x00\x00\x00\x00\x00\x00\x00"
			responseParts := strings.Split(string(response), sep)
			require.Len(t, responseParts, 2)
			head := string(responseParts[0])
			body := string(responseParts[1])
			assert.JSONEq(t, test.expectedHead, head)
			assert.Equal(t, test.expectedBody, body)
			assert.NoError(t, test.response.Close())
		})
	}
}

type readCloser struct {
	closed bool
	err    error
	reader *strings.Reader
}

func (r *readCloser) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *readCloser) Close() error {
	r.closed = true
	return r.err
}

func TestLambdaFunctionURLStreamingResponsePropogatesInnerClose(t *testing.T) {
	for _, test := range []struct {
		name   string
		closer *readCloser
		err    error
	}{
		{
			"closer no err",
			&readCloser{},
			nil,
		},
		{
			"closer with err",
			&readCloser{err: errors.New("yolo")},
			errors.New("yolo"),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := &LambdaFunctionURLStreamingResponse{Body: test.closer}
			assert.Equal(t, test.err, response.Close())
			assert.True(t, test.closer.closed)
		})
	}
}
