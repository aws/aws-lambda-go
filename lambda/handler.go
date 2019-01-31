// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda/handlertrace"
)

type Handler interface {
	Invoke(ctx context.Context, payload []byte) ([]byte, error)
}

type functionHandler struct {
	takesContext     bool
	requestEventType *reflect.Type
	originalFunc     reflect.Value
}

func (handler functionHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	trace := handlertrace.FromContext(ctx)

	var args []reflect.Value
	if handler.takesContext {
		args = append(args, reflect.ValueOf(ctx))
	}
	if nil != handler.requestEventType {
		event := reflect.New(*handler.requestEventType)
		if err := json.Unmarshal(payload, event.Interface()); err != nil {
			return nil, err
		}
		if nil != trace.RequestEvent {
			trace.RequestEvent(ctx, event.Elem().Interface())
		}
		args = append(args, event.Elem())
	}

	response := handler.originalFunc.Call(args)

	if len(response) > 0 {
		if err, ok := response[len(response)-1].Interface().(error); ok && err != nil {
			return nil, err
		}
	}
	var val interface{}
	if len(response) > 1 {
		val = response[0].Interface()

		if nil != trace.ResponseEvent {
			trace.ResponseEvent(ctx, val)
		}
	}
	return json.Marshal(val)
}

type errHandler struct{ e error }

func (e errHandler) Invoke(context.Context, []byte) ([]byte, error) { return nil, e.e }

func errorHandler(e error) Handler { return errHandler{e: e} }

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
	if handler.NumOut() > 2 {
		return fmt.Errorf("handler may not return more than two values")
	} else if handler.NumOut() > 1 {
		if !handler.Out(1).Implements(errorType) {
			return fmt.Errorf("handler returns two values, but the second does not implement error")
		}
	} else if handler.NumOut() == 1 {
		if !handler.Out(0).Implements(errorType) {
			return fmt.Errorf("handler returns a single value, but it does not implement error")
		}
	}
	return nil
}

// NewHandler creates a base lambda handler from the given handler function. The
// returned Handler performs JSON deserialization and deserialization, and
// delegates to the input handler function.  The handler function parameter must
// satisfy the rules documented by Start.  If handlerFunc is not a valid
// handler, the returned Handler simply reports the validation error.
func NewHandler(handlerFunc interface{}) Handler {
	if handlerFunc == nil {
		return errorHandler(fmt.Errorf("handler is nil"))
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

	var requestEventType *reflect.Type
	if (handlerType.NumIn() == 1 && !takesContext) || handlerType.NumIn() == 2 {
		eventType := handlerType.In(handlerType.NumIn() - 1)
		requestEventType = &eventType
	}

	return functionHandler{
		takesContext:     takesContext,
		requestEventType: requestEventType,
		originalFunc:     handler,
	}
}
