//go:build go1.18
// +build go1.18

// Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package lambdaurl

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/function-url-request-with-headers-and-cookies-and-text-body.json
var helloRequest []byte

//go:embed testdata/function-url-domain-only-get-request.json
var domainOnlyGetRequest []byte

//go:embed testdata/function-url-domain-only-get-request-trailing-slash.json
var domainOnlyWithSlashGetRequest []byte

//go:embed testdata/function-url-domain-only-request-with-base64-encoded-body.json
var base64EncodedBodyRequest []byte

func TestWrap(t *testing.T) {
	for name, params := range map[string]struct {
		input         []byte
		handler       http.HandlerFunc
		expectStatus  int
		expectBody    string
		expectHeaders map[string]string
		expectCookies []string
	}{
		"hello": {
			input: helloRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Hello", "world1")
				w.Header().Add("Hello", "world2")
				http.SetCookie(w, &http.Cookie{Name: "yummy", Value: "cookie"})
				http.SetCookie(w, &http.Cookie{Name: "yummy", Value: "cake"})
				http.SetCookie(w, &http.Cookie{Name: "fruit", Value: "banana", Expires: time.Date(2000, time.January, 0, 0, 0, 0, 0, time.UTC)})
				for _, c := range r.Cookies() {
					http.SetCookie(w, c)
				}

				w.WriteHeader(http.StatusTeapot)
				encoder := json.NewEncoder(w)
				_ = encoder.Encode(struct{ RequestQueryParams, Method any }{r.URL.Query(), r.Method})
			},
			expectStatus: http.StatusTeapot,
			expectHeaders: map[string]string{
				"Hello": "world1,world2",
			},
			expectCookies: []string{
				"yummy=cookie",
				"yummy=cake",
				"fruit=banana; Expires=Fri, 31 Dec 1999 00:00:00 GMT",
				"foo=bar",
				"hello=hello",
			},
			expectBody: `{"RequestQueryParams":{"foo":["bar"],"hello":["world"]},"Method":"POST"}` + "\n",
		},
		"mux": {
			input: helloRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				log.Println(r.URL)
				mux := http.NewServeMux()
				mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					_, _ = w.Write([]byte("Hello World!"))
				})
				mux.ServeHTTP(w, r)
			},
			expectStatus: 200,
			expectBody:   "Hello World!",
		},
		"get-implicit-trailing-slash": {
			input: domainOnlyGetRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				_ = encoder.Encode(r.Method)
				_ = encoder.Encode(r.URL.String())
			},
			expectStatus: http.StatusOK,
			expectBody:   "\"GET\"\n\"https://lambda-url-id.lambda-url.us-west-2.on.aws/\"\n",
		},
		"get-explicit-trailing-slash": {
			input: domainOnlyWithSlashGetRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				_ = encoder.Encode(r.Method)
				_ = encoder.Encode(r.URL.String())
			},
			expectStatus: http.StatusOK,
			expectBody:   "\"GET\"\n\"https://lambda-url-id.lambda-url.us-west-2.on.aws/\"\n",
		},
		"empty handler": {
			input:        helloRequest,
			handler:      func(w http.ResponseWriter, r *http.Request) {},
			expectStatus: http.StatusOK,
		},
		"base64request": {
			input: base64EncodedBodyRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.Copy(w, r.Body)
			},
			expectStatus: http.StatusOK,
			expectBody:   "<idk/>",
		},
	} {
		t.Run(name, func(t *testing.T) {
			handler := Wrap(params.handler)
			var req events.LambdaFunctionURLRequest
			require.NoError(t, json.Unmarshal(params.input, &req))
			res, err := handler(context.Background(), &req)
			require.NoError(t, err)
			resultBodyBytes, err := ioutil.ReadAll(res)
			require.NoError(t, err)
			resultHeaderBytes, resultBodyBytes, ok := bytes.Cut(resultBodyBytes, []byte{0, 0, 0, 0, 0, 0, 0, 0})
			require.True(t, ok)
			var resultHeader struct {
				StatusCode int
				Headers    map[string]string
				Cookies    []string
			}
			require.NoError(t, json.Unmarshal(resultHeaderBytes, &resultHeader))
			assert.Equal(t, params.expectBody, string(resultBodyBytes))
			assert.Equal(t, params.expectStatus, resultHeader.StatusCode)
			assert.Equal(t, params.expectHeaders, resultHeader.Headers)
			assert.Equal(t, params.expectCookies, resultHeader.Cookies)
		})
	}
}
