// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFatalErrors(t *testing.T) {
	ts, record := runtimeAPIServer(``, 100)
	defer ts.Close()
	handler := NewHandler(func() error {
		panic(errors.New("a fatal error"))
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedErrorMessage := "calling the handler function resulted in a panic, the process should exit"
	assert.EqualError(t, startRuntimeAPILoop(endpoint, handler), expectedErrorMessage)
	assert.Equal(t, 1, record.nGets)
	var invokeErr messages.InvokeResponse_Error
	err := json.Unmarshal(record.responses[0], &invokeErr)
	assert.NoError(t, err)
	assert.NotNil(t, invokeErr.StackTrace)
	assert.Equal(t, "errorString", invokeErr.Type)
	assert.Equal(t, "a fatal error", invokeErr.Message)
}

func TestRuntimeAPILoop(t *testing.T) {
	nInvokes := 10

	ts, record := runtimeAPIServer(``, nInvokes)
	defer ts.Close()

	n := 0
	handler := NewHandler(func(ctx context.Context) (string, error) {
		n += 1
		if n%3 == 0 {
			return "", errors.New("error time!")
		}
		return "Hello!", nil
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedError := fmt.Sprintf("failed to GET http://%s/2018-06-01/runtime/invocation/next: got unexpected status code: 410", endpoint)
	assert.EqualError(t, startRuntimeAPILoop(endpoint, handler), expectedError)
	assert.Equal(t, nInvokes+1, record.nGets)
	assert.Equal(t, nInvokes, record.nPosts)
}

func TestCustomErrorMarshaling(t *testing.T) {
	type CustomError struct{ error }
	errors := []error{
		errors.New("boring"),
		CustomError{errors.New("Something bad happened!")},
		messages.InvokeResponse_Error{Type: "yolo", Message: "hello"},
	}
	expected := []string{
		`{ "errorType": "errorString", "errorMessage": "boring"}`,
		`{ "errorType": "CustomError", "errorMessage": "Something bad happened!" }`,
		`{ "errorType": "yolo", "errorMessage": "hello" }`,
	}
	require.Equal(t, len(errors), len(expected))
	ts, record := runtimeAPIServer(``, len(errors))
	defer ts.Close()
	n := 0
	handler := NewHandler(func() error {
		defer func() { n++ }()
		return errors[n]
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedError := fmt.Sprintf("failed to GET http://%s/2018-06-01/runtime/invocation/next: got unexpected status code: 410", endpoint)
	assert.EqualError(t, startRuntimeAPILoop(endpoint, handler), expectedError)
	for i := range errors {
		assert.JSONEq(t, expected[i], string(record.responses[i]))
		assert.Equal(t, contentTypeJSON, record.contentTypes[i])
	}
}

func TestXRayCausePlumbing(t *testing.T) {
	errors := []error{
		errors.New("barf"),
		messages.InvokeResponse_Error{
			Type:    "yoloError",
			Message: "hello yolo",
			StackTrace: []*messages.InvokeResponse_Error_StackFrame{
				{Label: "yolo", Path: "yolo", Line: 2},
				{Label: "hi", Path: "hello/hello", Line: 12},
			},
		},
		messages.InvokeResponse_Error{
			Type:    "yoloError",
			Message: "hello yolo",
			StackTrace: []*messages.InvokeResponse_Error_StackFrame{
			},
		},
	}
	wd, _ := os.Getwd()
	expected := []string{
	    `{
		    "working_directory":"` + wd + `", 
		    "paths": [], 
		    "exceptions": [{ 
			"type": "errorString", 
			"message": "barf", 
			"stack": []
		    }]
		}`,
		`{
		    "working_directory":"` + wd + `", 
		    "paths": ["yolo", "hello/hello"], 
		    "exceptions": [{ 
			"type": "yoloError", 
			"message": "hello yolo", 
			"stack": [
			    {"label": "yolo", "path": "yolo", "line": 2},
			    {"label": "hi", "path": "hello/hello", "line": 12}
			]
		    }]
		}`,
		`{
		    "working_directory":"` + wd + `", 
		    "paths": [], 
		    "exceptions": [{ 
			"type": "yoloError", 
			"message": "hello yolo", 
			"stack": [
			]
		    }]
		}`,

	}
	require.Equal(t, len(errors), len(expected))
	ts, record := runtimeAPIServer(``, len(errors))
	defer ts.Close()
	n := 0
	handler := NewHandler(func() error {
		defer func() { n++ }()
		return errors[n]
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedError := fmt.Sprintf("failed to GET http://%s/2018-06-01/runtime/invocation/next: got unexpected status code: 410", endpoint)
	assert.EqualError(t, startRuntimeAPILoop(endpoint, handler), expectedError)
	for i := range errors {
		assert.JSONEq(t, expected[i], string(record.xrayCauses[i]))
	}

}

func TestRuntimeAPIContextPlumbing(t *testing.T) {
	handler := NewHandler(func(ctx context.Context) (interface{}, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		deadline, _ := ctx.Deadline()
		return struct {
			Context    *lambdacontext.LambdaContext
			TraceID    string
			EnvTraceID string
			Deadline   int64
		}{
			Context:    lc,
			TraceID:    ctx.Value("x-amzn-trace-id").(string),
			EnvTraceID: os.Getenv("_X_AMZN_TRACE_ID"),
			Deadline:   deadline.UnixNano() / nsPerMS,
		}, nil
	})

	ts, record := runtimeAPIServer(``, 1)
	defer ts.Close()

	endpoint := strings.Split(ts.URL, "://")[1]
	expectedError := fmt.Sprintf("failed to GET http://%s/2018-06-01/runtime/invocation/next: got unexpected status code: 410", endpoint)
	assert.EqualError(t, startRuntimeAPILoop(endpoint, handler), expectedError)

	expected := `
	{
		"Context": {
			"AwsRequestID": "dummyid",
			"InvokedFunctionArn": "dummyarn",
			"Identity": {
				"CognitoIdentityID": "dummyident",
				"CognitoIdentityPoolID": "dummypool"
			},
			"ClientContext": {
				"Client": {
					"installation_id": "dummyinstallid",
					"app_title": "dummytitle",
					"app_version_code": "dummycode",
					"app_package_name": "dummyname"
				},
				"env": null,
				"custom": null
			}
		},
		"TraceID": "its-xray-time",
		"EnvTraceID": "its-xray-time",
		"Deadline": 22
	}
	`
	assert.JSONEq(t, expected, string(record.responses[0]))
}

func TestReadPayload(t *testing.T) {
	ts, record := runtimeAPIServer(`{"message": "I am craving tacos"}`, 1)
	defer ts.Close()

	handler := NewHandler(func(event struct{ Message string }) (string, error) {
		length := utf8.RuneCountInString(event.Message)
		reversed := make([]rune, length)
		for i, v := range event.Message {
			reversed[length-i-1] = v
		}
		return string(reversed), nil
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	_ = startRuntimeAPILoop(endpoint, handler)
	assert.Equal(t, `"socat gnivarc ma I"`, string(record.responses[0]))
	assert.Equal(t, contentTypeJSON, record.contentTypes[0])
}

type readCloser struct {
	closed bool
	reader *strings.Reader
}

func (r *readCloser) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *readCloser) Close() error {
	r.closed = true
	return nil
}

func TestBinaryResponseDefaultContentType(t *testing.T) {
	ts, record := runtimeAPIServer(`{"message": "I am craving tacos"}`, 1)
	defer ts.Close()

	handler := NewHandler(func(event struct{ Message string }) (io.Reader, error) {
		length := utf8.RuneCountInString(event.Message)
		reversed := make([]rune, length)
		for i, v := range event.Message {
			reversed[length-i-1] = v
		}
		return strings.NewReader(string(reversed)), nil
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	_ = startRuntimeAPILoop(endpoint, handler)
	assert.Equal(t, `socat gnivarc ma I`, string(record.responses[0]))
	assert.Equal(t, contentTypeBytes, record.contentTypes[0])
}

func TestBinaryResponseDoesNotLeakResources(t *testing.T) {
	numResponses := 3
	responses := make([]*readCloser, numResponses)
	for i := 0; i < numResponses; i++ {
		responses[i] = &readCloser{closed: false, reader: strings.NewReader(fmt.Sprintf("hello %d", i))}
	}
	timesCalled := 0
	handler := NewHandler(func() (res io.Reader, _ error) {
		res = responses[timesCalled]
		timesCalled++
		return
	})

	ts, record := runtimeAPIServer(`{}`, numResponses)
	defer ts.Close()
	endpoint := strings.Split(ts.URL, "://")[1]
	_ = startRuntimeAPILoop(endpoint, handler)

	for i := 0; i < numResponses; i++ {
		assert.Equal(t, contentTypeBytes, record.contentTypes[i])
		assert.Equal(t, fmt.Sprintf("hello %d", i), string(record.responses[i]))
		assert.True(t, responses[i].closed)
	}
}

func TestContextDeserializationErrors(t *testing.T) {
	badClientContext := defaultInvokeMetadata()
	badClientContext.clientContext = `{ not json }`

	badCognito := defaultInvokeMetadata()
	badCognito.cognito = `{ not json }`

	badDeadline := defaultInvokeMetadata()
	badDeadline.deadline = `yolo`

	badMetadata := []eventMetadata{badClientContext, badCognito, badDeadline}

	ts, record := runtimeAPIServer(`{}`, len(badMetadata), badMetadata...)
	defer ts.Close()
	handler := NewHandler(func(ctx context.Context) (*lambdacontext.LambdaContext, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		return lc, nil
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	_ = startRuntimeAPILoop(endpoint, handler)

	assert.JSONEq(t, `{
	    "errorMessage":"failed to unmarshal client context json: invalid character 'n' looking for beginning of object key string",
	    "errorType":"errorString"
	}`, string(record.responses[0]))

	assert.JSONEq(t, `{
	    "errorMessage":"failed to unmarshal cognito identity json: invalid character 'n' looking for beginning of object key string",
	    "errorType":"errorString"
	}`, string(record.responses[1]))

	assert.JSONEq(t, `{
	    "errorMessage":"failed to parse deadline: strconv.ParseInt: parsing \"yolo\": invalid syntax",
	    "errorType":"errorString"
	}`, string(record.responses[2]))
}

type invalidPayload struct{}

func (invalidPayload) MarshalJSON() ([]byte, error) {
	return nil, errors.New(`some error that contains '"'`)
}

func TestSafeMarshal_SerializationError(t *testing.T) {
	payload := safeMarshal(invalidPayload{})
	want := `{"errorMessage":"json: error calling MarshalJSON for type lambda.invalidPayload: some error that contains '\"'","errorType":"Runtime.SerializationError"}`
	assert.Equal(t, want, string(payload))
}

type requestRecord struct {
	nGets        int
	nPosts       int
	responses    [][]byte
	contentTypes []string
	xrayCauses   []string
}

type eventMetadata struct {
	clientContext string
	cognito       string
	xray          string
	deadline      string
	requestID     string
	functionARN   string
}

func defaultInvokeMetadata() eventMetadata {
	return eventMetadata{
		clientContext: `{
			"Client": {
				"app_title": "dummytitle",
				"installation_id": "dummyinstallid",
				"app_version_code": "dummycode",
				"app_package_name": "dummyname"
			}
		}`,
		cognito: `{
			"cognitoIdentityId": "dummyident", 
			"cognitoIdentityPoolId": "dummypool"
		}`,
		xray:        "its-xray-time",
		requestID:   "dummyid",
		deadline:    "22",
		functionARN: "dummyarn",
	}
}

func runtimeAPIServer(eventPayload string, failAfter int, overrides ...eventMetadata) (*httptest.Server, *requestRecord) {
	numInvokesRequested := 0
	record := &requestRecord{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			metadata := defaultInvokeMetadata()
			if numInvokesRequested < len(overrides) {
				metadata = overrides[numInvokesRequested]
			}
			record.nGets++
			numInvokesRequested++
			if numInvokesRequested > failAfter {
				w.WriteHeader(http.StatusGone)
				_, _ = w.Write([]byte("END THE TEST!"))
			}
			w.Header().Add(string(headerAWSRequestID), metadata.requestID)
			w.Header().Add(string(headerDeadlineMS), metadata.deadline)
			w.Header().Add(string(headerInvokedFunctionARN), metadata.functionARN)
			w.Header().Add(string(headerClientContext), metadata.clientContext)
			w.Header().Add(string(headerCognitoIdentity), metadata.cognito)
			w.Header().Add(string(headerTraceID), metadata.xray)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(eventPayload))
		case http.MethodPost:
			record.nPosts++
			response := bytes.NewBuffer(nil)
			_, _ = io.Copy(response, r.Body)
			_ = r.Body.Close()
			w.WriteHeader(http.StatusAccepted)
			record.responses = append(record.responses, response.Bytes())
			record.contentTypes = append(record.contentTypes, r.Header.Get("Content-Type"))
			record.xrayCauses = append(record.xrayCauses, r.Header.Get(headerXRayErrorCause))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	return ts, record
}
