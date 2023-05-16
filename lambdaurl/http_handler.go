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

type options struct {
	detectContentType bool
	bufferHeader      bool
}

type Option func(*options)

// WithDetectContentType sets the behavior of content type detection when the Content-Type header is not already provided.
// When true, the first Write call will pass the intial bytes to http.DetectContentType.
// When false, and if no Content-Type is provided, no Content-Type will be sent back to Lambda,
// and the Lambda Function URL will fallback to it's default.
//
// Note: The http.ResponseWriter passed to the handler is unbuffered.
// This may result in different Content-Type headers in the Function URL response when compared to http.ListenAndServe.
func WithDetectContentType(detectContentType bool) func(*options) {
	return func(options *options) {
		options.detectContentType = detectContentType
	}
}

type httpResponseWriter struct {
	options options
	header  http.Header
	writer  io.Writer
	once    sync.Once
	ready   chan<- header
}

type header struct {
	code   int
	header http.Header
}

func (w *httpResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *httpResponseWriter) Write(p []byte) (int, error) {
	w.writeHeader(http.StatusOK, p)
	return w.writer.Write(p)
}

func (w *httpResponseWriter) WriteHeader(statusCode int) {
	w.writeHeader(statusCode, nil)
}

func (w *httpResponseWriter) writeHeader(statusCode int, initialPayload []byte) {
	w.once.Do(func() {
		if w.options.detectContentType {
			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", http.DetectContentType(initialPayload))
			}
		}
		w.ready <- header{code: statusCode, header: w.header}
	})
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
func Wrap(handler http.Handler) func(context.Context, *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {
	return wrap(handler)
}

// WrapWithOptions converts an http.Handler into a Lambda request handler.
//
// Only Lambda Function URLs configured with `InvokeMode: RESPONSE_STREAM` are supported with the returned handler.
// The response body of the handler will conform to the content-type `application/vnd.awslambda.http-integration-response`.
func WrapWithOptions(handler http.Handler, options ...Option) func(context.Context, *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {
	return wrap(handler, options...)
}

func wrap(handler http.Handler, options ...Option) func(context.Context, *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {
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
		ready := make(chan header) // Signals when it's OK to start returning the response body to Lambda
		r, w := io.Pipe()
		go func() {
			defer close(ready)
			defer w.Close() // TODO: recover and CloseWithError the any panic value once the runtime API client supports plumbing fatal errors through the reader
			responseWriter := &httpResponseWriter{writer: w, ready: ready}
			defer responseWriter.Write(nil) // force default status, headers, content type detection, if none occured during the execution of the handler
			for _, f := range options {
				f(&responseWriter.options)
			}
			handler.ServeHTTP(responseWriter, httpRequest)
		}()
		header := <-ready
		response := &events.LambdaFunctionURLStreamingResponse{
			Body:       r,
			StatusCode: header.code,
		}
		if len(header.header) > 0 {
			response.Headers = make(map[string]string, len(header.header))
			for k, v := range header.header {
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
