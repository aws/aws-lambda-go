// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda/handlertrace"
)

type Handler interface {
	Invoke(ctx context.Context, payload []byte) ([]byte, error)
}

type handlerOptions struct {
	Handler
	baseContext              context.Context
	jsonResponseEscapeHTML   bool
	jsonResponseIndentPrefix string
	jsonResponseIndentValue  string
	enableSIGTERM            bool
	sigtermCallbacks         []func()
}

type Option func(*handlerOptions)

// WithContext is a HandlerOption that sets the base context for all invocations of the handler.
//
// Usage:
//
//	lambda.StartWithOptions(
//	 	func (ctx context.Context) (string, error) {
//	 		return ctx.Value("foo"), nil
//	 	},
//	 	lambda.WithContext(context.WithValue(context.Background(), "foo", "bar"))
//	)
func WithContext(ctx context.Context) Option {
	return Option(func(h *handlerOptions) {
		h.baseContext = ctx
	})
}

// WithSetEscapeHTML sets the SetEscapeHTML argument on the underlying json encoder
//
// Usage:
//
//	lambda.StartWithOptions(
//		func () (string, error) {
//			return "<html><body>hello!></body></html>", nil
//		},
//		lambda.WithSetEscapeHTML(true),
//	)
func WithSetEscapeHTML(escapeHTML bool) Option {
	return Option(func(h *handlerOptions) {
		h.jsonResponseEscapeHTML = escapeHTML
	})
}

// WithSetIndent sets the SetIndent argument on the underling json encoder
//
// Usage:
//
//	lambda.StartWithOptions(
//		func (event any) (any, error) {
//			return event, nil
//		},
//		lambda.WithSetIndent(">"," "),
//	)
func WithSetIndent(prefix, indent string) Option {
	return Option(func(h *handlerOptions) {
		h.jsonResponseIndentPrefix = prefix
		h.jsonResponseIndentValue = indent
	})
}

// WithEnableSIGTERM enables SIGTERM behavior within the Lambda platform on container spindown.
// SIGKILL will occur ~500ms after SIGTERM.
// Optionally, an array of callback functions to run on SIGTERM may be provided.
//
// Usage:
//
//	lambda.StartWithOptions(
//	    func (event any) (any, error) {
//			return event, nil
//		},
//		lambda.WithEnableSIGTERM(func() {
//			log.Print("function container shutting down...")
//		})
//	)
func WithEnableSIGTERM(callbacks ...func()) Option {
	return Option(func(h *handlerOptions) {
		h.sigtermCallbacks = append(h.sigtermCallbacks, callbacks...)
		h.enableSIGTERM = true
	})
}

func validateArguments(handler reflect.Type) (bool, error) {
	handlerTakesContext := false
	if handler.NumIn() > 2 {
		return false, fmt.Errorf("handlers may not take more than two arguments, but handler takes %d", handler.NumIn())
	} else if handler.NumIn() > 0 {
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
		argumentType := handler.In(0)
		handlerTakesContext = argumentType.Implements(contextType)
		if handler.NumIn() > 1 && !handlerTakesContext {
			return false, fmt.Errorf("handler takes two arguments, but the first is not Context. got %s", argumentType.Kind())
		}
	}

	return handlerTakesContext, nil
}

func validateReturns(handler reflect.Type) error {
	errorType := reflect.TypeOf((*error)(nil)).Elem()

	switch n := handler.NumOut(); {
	case n > 2:
		return fmt.Errorf("handler may not return more than two values")
	case n > 1:
		if !handler.Out(1).Implements(errorType) {
			return fmt.Errorf("handler returns two values, but the second does not implement error")
		}
	case n == 1:
		if !handler.Out(0).Implements(errorType) {
			return fmt.Errorf("handler returns a single value, but it does not implement error")
		}
	}

	return nil
}

// NewHandler creates a base lambda handler from the given handler function. The
// returned Handler performs JSON serialization and deserialization, and
// delegates to the input handler function. The handler function parameter must
// satisfy the rules documented by Start. If handlerFunc is not a valid
// handler, the returned Handler simply reports the validation error.
func NewHandler(handlerFunc interface{}) Handler {
	return NewHandlerWithOptions(handlerFunc)
}

// NewHandlerWithOptions creates a base lambda handler from the given handler function. The
// returned Handler performs JSON serialization and deserialization, and
// delegates to the input handler function. The handler function parameter must
// satisfy the rules documented by Start. If handlerFunc is not a valid
// handler, the returned Handler simply reports the validation error.
func NewHandlerWithOptions(handlerFunc interface{}, options ...Option) Handler {
	return newHandler(handlerFunc, options...)
}

func newHandler(handlerFunc interface{}, options ...Option) *handlerOptions {
	if h, ok := handlerFunc.(*handlerOptions); ok {
		return h
	}
	h := &handlerOptions{
		baseContext:              context.Background(),
		jsonResponseEscapeHTML:   false,
		jsonResponseIndentPrefix: "",
		jsonResponseIndentValue:  "",
	}
	for _, option := range options {
		option(h)
	}
	if h.enableSIGTERM {
		enableSIGTERM(h.sigtermCallbacks)
	}
	h.Handler = reflectHandler(handlerFunc, h)
	return h
}

type bytesHandlerFunc func(context.Context, []byte) ([]byte, error)

func (h bytesHandlerFunc) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return h(ctx, payload)
}
func errorHandler(err error) Handler {
	return bytesHandlerFunc(func(_ context.Context, _ []byte) ([]byte, error) {
		return nil, err
	})
}

func reflectHandler(handlerFunc interface{}, h *handlerOptions) Handler {
	if handlerFunc == nil {
		return errorHandler(errors.New("handler is nil"))
	}

	if handler, ok := handlerFunc.(Handler); ok {
		return handler
	}

	handler := reflect.ValueOf(handlerFunc)
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		return errorHandler(fmt.Errorf("handler kind %s is not %s", handlerType.Kind(), reflect.Func))
	}

	takesContext, err := validateArguments(handlerType)
	if err != nil {
		return errorHandler(err)
	}

	if err := validateReturns(handlerType); err != nil {
		return errorHandler(err)
	}

	return bytesHandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
		in := bytes.NewBuffer(payload)
		out := bytes.NewBuffer(nil)
		decoder := json.NewDecoder(in)
		encoder := json.NewEncoder(out)
		encoder.SetEscapeHTML(h.jsonResponseEscapeHTML)
		encoder.SetIndent(h.jsonResponseIndentPrefix, h.jsonResponseIndentValue)

		trace := handlertrace.FromContext(ctx)

		// construct arguments
		var args []reflect.Value
		if takesContext {
			args = append(args, reflect.ValueOf(ctx))
		}
		if (handlerType.NumIn() == 1 && !takesContext) || handlerType.NumIn() == 2 {
			eventType := handlerType.In(handlerType.NumIn() - 1)
			event := reflect.New(eventType)
			if err := decoder.Decode(event.Interface()); err != nil {
				return nil, err
			}
			if nil != trace.RequestEvent {
				trace.RequestEvent(ctx, event.Elem().Interface())
			}
			args = append(args, event.Elem())
		}

		response := handler.Call(args)

		// return the error, if any
		if len(response) > 0 {
			if errVal, ok := response[len(response)-1].Interface().(error); ok && errVal != nil {
				return nil, errVal
			}
		}
		// set the response value, if any
		var val interface{}
		if len(response) > 1 {
			val = response[0].Interface()
			if nil != trace.ResponseEvent {
				trace.ResponseEvent(ctx, val)
			}
		}
		if err := encoder.Encode(val); err != nil {
			return nil, err
		}

		responseBytes := out.Bytes()
		// back-compat, strip the encoder's trailing newline unless WithSetIndent was used
		if h.jsonResponseIndentValue == "" && h.jsonResponseIndentPrefix == "" {
			return responseBytes[:len(responseBytes)-1], nil
		}

		return responseBytes, nil
	})
}
