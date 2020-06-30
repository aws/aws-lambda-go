// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestFatalErrors(t *testing.T) {
	ts, record := runtimeAPIServer(``, 100)
	defer ts.Close()
	handler := NewHandler(func() error {
		panic(errors.New("a fatal error"))
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedErrorMessage := "calling the handler function resulted in a panic, the process should exit"
	assert.EqualError(t, startRuntimeAPILoop(context.Background(), endpoint, handler), expectedErrorMessage)
	assert.Equal(t, 1, record.nGets)
	assert.Equal(t, 1, record.nGets)
}

func TestRuntimeAPILoop(t *testing.T) {
	nInvokes := 10

	ts, record := runtimeAPIServer(``, nInvokes)
	defer ts.Close()

	n := 0
	handler := NewHandler(func() (string, error) {
		n += 1
		if n%3 == 0 {
			return "", errors.New("error time!")
		}
		return "Hello!", nil
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedError := fmt.Sprintf("failed to GET http://%s/2018-06-01/runtime/invocation/next: got unexpected status code: 410", endpoint)
	assert.EqualError(t, startRuntimeAPILoop(context.Background(), endpoint, handler), expectedError)
	assert.Equal(t, nInvokes+1, record.nGets)
	assert.Equal(t, nInvokes, record.nPosts)
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
	_ = startRuntimeAPILoop(context.Background(), endpoint, handler)
	assert.Equal(t, `"socat gnivarc ma I"`, string(record.responses[0]))

}

type requestRecord struct {
	nGets     int
	nPosts    int
	responses [][]byte
}

func runtimeAPIServer(eventPayload string, failAfter int) (*httptest.Server, *requestRecord) {
	numInvokesRequested := 0
	record := &requestRecord{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			record.nGets++
			numInvokesRequested++
			if numInvokesRequested > failAfter {
				w.WriteHeader(http.StatusGone)
				_, _ = w.Write([]byte("END THE TEST!"))
			}
			w.Header().Add(string(headerAWSRequestID), "dummy-request-id")
			w.Header().Add(string(headerDeadlineMS), "22")
			w.Header().Add(string(headerInvokedFunctionARN), "anarn")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(eventPayload))
		case http.MethodPost:
			record.nPosts++
			response := bytes.NewBuffer(nil)
			_, _ = io.Copy(response, r.Body)
			r.Body.Close()
			w.WriteHeader(http.StatusAccepted)
			record.responses = append(record.responses, response.Bytes())
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	return ts, record
}
