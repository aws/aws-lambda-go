//go:build go1.18
// +build go1.18

// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
)

// HandlerFunc represents a valid input as described by Start
type HandlerFunc[TIn, TOut any] interface {
	func() |
		func(TIn) |
		func() error |
		func(TIn) error |
		func() (TOut, error) |
		func(TIn) (TOut, error) |
		func(context.Context, TIn) |
		func(context.Context, TIn) error |
		func(context.Context, TIn) (TOut, error)
}

// StartWithOptionsTypeSafe is the same as StartWithOptions except that it takes a generic input
// so that the function signature can be validated at compile time.
// The caller can supply "any" for TIn or TOut if the input function does not use that argument or return value.
//
// Examples:
//
// TIn and TOut ignored
//
//	StartWithOptionsTypeSafe[any, any](func() {
//		fmt.Println("Hello world")
//	})
//
// TIn used and TOut ignored
//
//	type event events.APIGatewayV2HTTPRequest
//	StartWithOptionsTypeSafe[event, any](func(e event) {
//		fmt.Printf("Version: %s", e.Version)
//	})
//
// TIn ignored, TOut used and an error returned
//
//	type resp events.APIGatewayV2HTTPResponse
//	StartWithOptionsTypeSafe[any, resp](func() (resp, error) {
//		return resp{Body: "hello, world"}, nil
//	})
func StartWithOptionsTypeSafe[TIn any, TOut any, H HandlerFunc[TIn, TOut]](handler H, options ...Option) {
	start(newHandler(handler, options...))
}
