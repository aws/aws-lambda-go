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
	"os"
	"os/exec"
	"path"
	"strings"
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
		input             []byte
		handler           http.HandlerFunc
		detectContentType bool
		expectStatus      int
		expectBody        string
		expectHeaders     map[string]string
		expectCookies     []string
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
			expectStatus:  http.StatusTeapot,
			expectHeaders: map[string]string{"Hello": "world1,world2"},
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
		"write status code only": {
			input: helloRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			},
			expectStatus: http.StatusAccepted,
		},
		"base64request": {
			input: base64EncodedBodyRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.Copy(w, r.Body)
			},
			expectStatus: http.StatusOK,
			expectBody:   "<idk/>",
		},
		"detect content type: write status code only": {
			input: helloRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			},
			detectContentType: true,
			expectStatus:      http.StatusAccepted,
			expectHeaders: map[string]string{
				"Content-Type": "application/octet-stream",
			},
		},
		"detect content type: empty handler": {
			input: helloRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			detectContentType: true,
			expectStatus:      http.StatusOK,
			expectHeaders: map[string]string{
				"Content-Type": "application/octet-stream",
			},
		},
		"detect content type: writes html": {
			input: helloRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("<!DOCTYPE HTML><html></html>"))
			},
			detectContentType: true,
			expectBody:        "<!DOCTYPE HTML><html></html>",
			expectStatus:      http.StatusOK,
			expectHeaders: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		"detect content type: writes zeros": {
			input: helloRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte{0, 0, 0, 0, 0})
			},
			detectContentType: true,
			expectBody:        "\x00\x00\x00\x00\x00",
			expectStatus:      http.StatusOK,
			expectHeaders: map[string]string{
				"Content-Type": "application/octet-stream",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			handler := Wrap(params.handler)
			var req events.LambdaFunctionURLRequest
			require.NoError(t, json.Unmarshal(params.input, &req))
			ctx := context.WithValue(context.Background(), detectContentTypeContextKey{}, params.detectContentType)
			res, err := handler(ctx, &req)
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

func TestRequestContext(t *testing.T) {
	var req *events.LambdaFunctionURLRequest
	require.NoError(t, json.Unmarshal(helloRequest, &req))
	handler := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqFromContext, exists := RequestFromContext(r.Context())
		require.True(t, exists)
		require.NotNil(t, reqFromContext)
		assert.Equal(t, req, reqFromContext)
	}))
	_, err := handler(context.Background(), req)
	require.NoError(t, err)
}

func TestStartViaEmulator(t *testing.T) {
	rieInvokeAPI := "http://localhost:8080/2015-03-31/functions/function/invocations"
	if _, err := exec.LookPath("aws-lambda-rie"); err != nil {
		t.Skipf("%v - install from https://github.com/aws/aws-lambda-runtime-interface-emulator/", err)
	}

	// compile our handler, it'll always run to timeout ensuring the SIGTERM is triggered by aws-lambda-rie
	testDir := t.TempDir()
	handlerBuild := exec.Command("go", "build", "-o", path.Join(testDir, "lambdaurl.handler"), "./testdata/lambdaurl.go")
	handlerBuild.Stderr = os.Stderr
	handlerBuild.Stdout = os.Stderr
	require.NoError(t, handlerBuild.Run())

	// run the runtime interface emulator, capture the logs for assertion
	cmd := exec.Command("aws-lambda-rie", "lambdaurl.handler")
	cmd.Env = []string{
		"PATH=" + testDir,
		"AWS_LAMBDA_FUNCTION_TIMEOUT=2",
	}
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	var logs string
	done := make(chan interface{}) // closed on completion of log flush
	go func() {
		logBytes, err := ioutil.ReadAll(stdout)
		require.NoError(t, err)
		logs = string(logBytes)
		close(done)
	}()
	require.NoError(t, cmd.Start())
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	// give a moment for the port to bind
	time.Sleep(500 * time.Millisecond)

	client := &http.Client{Timeout: 5 * time.Second} // http client timeout to prevent case from hanging on aws-lambda-rie
	resp, err := client.Post(rieInvokeAPI, "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	expected := "{\"statusCode\":200,\"headers\":{\"Content-Type\":\"text/html; charset=utf-8\"}}\x00\x00\x00\x00\x00\x00\x00\x00<!DOCTYPE HTML>\n<html>\n<body>\nHello World!\n</body>\n</html>\n"
	assert.Equal(t, expected, string(body))

	require.NoError(t, cmd.Process.Kill()) // now ensure the logs are drained
	<-done
	t.Logf("stdout:\n%s", logs)
}
