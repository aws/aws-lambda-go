// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientNext(t *testing.T) {
	dummyRequestID := "dummy-request-id"
	dummyPayload := `{"hello": "world"}`

	returnsBody := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/2018-06-01/runtime/invocation/next" {
			w.WriteHeader(http.StatusNotImplemented)
		}
		w.Header().Add(headerAWSRequestID, dummyRequestID)
		_, _ = w.Write([]byte(dummyPayload))
	}))
	defer returnsBody.Close()

	returnsNoBody := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/2018-06-01/runtime/invocation/next" {
			w.WriteHeader(http.StatusNotImplemented)
		}
		w.Header().Add(headerAWSRequestID, dummyRequestID)
		w.WriteHeader(http.StatusOK)
	}))
	defer returnsNoBody.Close()

	t.Run("handles regular response", func(t *testing.T) {
		invoke, err := newRuntimeAPIClient(serverAddress(returnsBody)).next()
		require.NoError(t, err)
		assert.Equal(t, dummyRequestID, invoke.id)
		assert.Equal(t, dummyPayload, string(invoke.payload))
	})

	t.Run("handles no body", func(t *testing.T) {
		invoke, err := newRuntimeAPIClient(serverAddress(returnsNoBody)).next()
		require.NoError(t, err)
		assert.Equal(t, dummyRequestID, invoke.id)
		assert.Equal(t, 0, len(invoke.payload))
	})
}

func TestClientDoneAndError(t *testing.T) {
	invokeID := "theid"

	var capturedErrors [][]byte
	var capturedResponses [][]byte

	acceptsResponses := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Logf("unexpected method: %s", r.Method)
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		if r.URL.Path != fmt.Sprintf("/2018-06-01/runtime/invocation/%s/error", invokeID) && r.URL.Path != fmt.Sprintf("/2018-06-01/runtime/invocation/%s/response", invokeID) {
			t.Logf("unexpected url path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		if strings.HasSuffix(r.URL.Path, "/error") {
			capturedErrors = append(capturedErrors, body)
		} else if strings.HasSuffix(r.URL.Path, "/response") {
			capturedResponses = append(capturedErrors, body)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer acceptsResponses.Close()

	client := newRuntimeAPIClient(serverAddress(acceptsResponses))
	inputPayloads := [][]byte{nil, {}, []byte("hello")}
	expectedPayloadsRecived := [][]byte{{}, {}, []byte("hello")} // nil payload expected to be read as empty bytes by the server
	for i, payload := range inputPayloads {
		invoke := &invoke{
			id:     invokeID,
			client: client,
		}
		t.Run(fmt.Sprintf("happy Done with payload[%d]", i), func(t *testing.T) {
			err := invoke.success(payload, contentTypeJSON)
			assert.NoError(t, err)
		})
		t.Run(fmt.Sprintf("happy Error with payload[%d]", i), func(t *testing.T) {
			err := invoke.failure(payload, contentTypeJSON)
			assert.NoError(t, err)
		})
	}
	assert.Equal(t, expectedPayloadsRecived, capturedErrors)
	assert.Equal(t, expectedPayloadsRecived, capturedResponses)
}

func TestInvalidRequestsForMalformedEndpoint(t *testing.T) {
	_, err := newRuntimeAPIClient("🚨").next()
	require.Error(t, err)
	err = (&invoke{client: newRuntimeAPIClient("🚨")}).success(nil, "")
	require.Error(t, err)
	err = (&invoke{client: newRuntimeAPIClient("🚨")}).failure(nil, "")
	require.Error(t, err)
}

func TestStatusCodes(t *testing.T) {
	for i := 200; i < 600; i++ {
		t.Run(fmt.Sprintf("status: %d", i), func(t *testing.T) {
			url := fmt.Sprintf("status-%d", i)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = ioutil.ReadAll(r.Body)
				w.WriteHeader(i)
			}))

			defer ts.Close()

			client := newRuntimeAPIClient(serverAddress(ts))
			invoke := &invoke{id: url, client: client}
			if i == http.StatusOK {
				t.Run("next should not error", func(t *testing.T) {
					_, err := client.next()
					require.NoError(t, err)
				})
			} else {
				t.Run("next should error", func(t *testing.T) {
					_, err := client.next()
					require.Error(t, err)
					if i != 301 && i != 302 && i != 303 {
						assert.Contains(t, err.Error(), "unexpected status code")
						assert.Contains(t, err.Error(), fmt.Sprintf("%d", i))
					}
				})
			}

			if i == http.StatusAccepted {
				t.Run("success should not error", func(t *testing.T) {
					err := invoke.success(nil, "")
					require.NoError(t, err)
				})
				t.Run("failure should not error", func(t *testing.T) {
					err := invoke.failure(nil, "")
					require.NoError(t, err)
				})
			} else {
				t.Run("success should error", func(t *testing.T) {
					err := invoke.success(nil, "")
					require.Error(t, err)
					if i != 301 && i != 302 && i != 303 {
						assert.Contains(t, err.Error(), "unexpected status code")
						assert.Contains(t, err.Error(), fmt.Sprintf("%d", i))
					}
				})
				t.Run("failure should error", func(t *testing.T) {
					err := invoke.failure(nil, "")
					require.Error(t, err)
					if i != 301 && i != 302 && i != 303 {
						assert.Contains(t, err.Error(), "unexpected status code")
						assert.Contains(t, err.Error(), fmt.Sprintf("%d", i))
					}
				})
			}
		})
	}
}

func serverAddress(ts *httptest.Server) string {
	return strings.Split(ts.URL, "://")[1]
}
