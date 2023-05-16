//go:build go1.18
// +build go1.18

// Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.

// Package lambdaurl serves requests from Lambda Function URLs using http.Handler.
package lambdaurl

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type httpResponseWriter struct {
	header http.Header
	writer io.Writer
	once   sync.Once
	status chan<- int
}

func (w *httpResponseWriter) Header() http.Header {
	return w.header
}

func (w *httpResponseWriter) Write(p []byte) (int, error) {
	w.once.Do(func() {
		w.detectContentType(p)
		w.status <- http.StatusOK
	})
	return w.writer.Write(p)
}

func (w *httpResponseWriter) WriteHeader(statusCode int) {
	w.once.Do(func() {
		w.detectContentType(nil)
		w.status <- statusCode
	})
}

func (w *httpResponseWriter) detectContentType(p []byte) {
	if w.header.Get("Content-Type") == "" {
		w.header.Set("Content-Type", http.DetectContentType(p))
	}
}

type requestContextKey struct{}

// RequestFromContext returns the *events.LambdaFunctionURLRequest from a context.
func RequestFromContext(ctx context.Context) (*events.LambdaFunctionURLRequest, bool) {
	req, ok := ctx.Value(requestContextKey{}).(*events.LambdaFunctionURLRequest)
	return req, ok
}

// Wrap converts an http.Handler into a Lambda request handler.
//
// Only Lambda Function URLs configured with `InvokeMode: RESPONSE_STREAM` are supported with the returned handler.
// The response body of the handler will conform to the content-type `application/vnd.awslambda.http-integration-response`.
//
// Note: The http.ResponseWriter passed to the handler is unbuffered.
// This may result in different Content-Type and Content-Length headers in the Function URL response when compared to http.ListenAndServe.
func Wrap(handler http.Handler) func(context.Context, *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {
	return func(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {
		var body io.Reader = strings.NewReader(request.Body)
		if request.IsBase64Encoded {
			body = base64.NewDecoder(base64.StdEncoding, body)
		}
		url := "https://" + request.RequestContext.DomainName + request.RawPath
		if request.RawQueryString != "" {
			url += "?" + request.RawQueryString
		}
		ctx = context.WithValue(ctx, requestContextKey{}, request)
		httpRequest, err := http.NewRequestWithContext(ctx, request.RequestContext.HTTP.Method, url, body)
		if err != nil {
			return nil, err
		}
		for k, v := range request.Headers {
			httpRequest.Header.Add(k, v)
		}
		status := make(chan int) // Signals when it's OK to start returning the response body to Lambda
		header := http.Header{}
		r, w := io.Pipe()
		go func() {
			defer close(status)
			defer w.Close() // TODO: recover and CloseWithError the any panic value once the runtime API client supports plumbing fatal errors through the reader
			responseWriter := &httpResponseWriter{writer: w, header: header, status: status}
			defer responseWriter.Write(nil)
			handler.ServeHTTP(responseWriter, httpRequest)
		}()
		response := &events.LambdaFunctionURLStreamingResponse{
			Body:       r,
			StatusCode: <-status,
		}
		if len(header) > 0 {
			response.Headers = make(map[string]string, len(header))
			for k, v := range header {
				if k == "Set-Cookie" {
					response.Cookies = v
				} else {
					response.Headers[k] = strings.Join(v, ",")
				}
			}
		}
		return response, nil
	}
}

// Start wraps a http.Handler and calls lambda.StartHandlerFunc
// Only supports:
//   - Lambda Function URLs configured with `InvokeMode: RESPONSE_STREAM`
//   - Lambda Functions using the `provided` or `provided.al2` runtimes.
//   - Lambda Functions using the `go1.x` runtime when compiled with `-tags lambda.norpc`
func Start(handler http.Handler, options ...lambda.Option) {
	lambda.StartHandlerFunc(Wrap(handler), options...)
}
