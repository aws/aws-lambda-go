//go:build go1.18
// +build go1.18

// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartWithOptionsTypeSafe(t *testing.T) {
	testCases := []struct {
		name         string
		handler      any
		takesContext bool
	}{
		{
			name:         "0 arg, 0 returns",
			handler:      func() {},
			takesContext: false,
		},
		{
			name:         "0 arg, 1 returns",
			handler:      func() error { return nil },
			takesContext: false,
		},
		{
			name:         "1 arg, 0 returns",
			handler:      func(any) {},
			takesContext: false,
		},
		{
			name:         "1 arg, 1 returns",
			handler:      func(any) error { return nil },
			takesContext: false,
		},
		{
			name:         "0 arg, 2 returns",
			handler:      func() (any, error) { return 1, nil },
			takesContext: false,
		},
		{
			name:         "1 arg, 2 returns",
			handler:      func(any) (any, error) { return 1, nil },
			takesContext: false,
		},
		{
			name:         "2 arg, 0 returns",
			handler:      func(context.Context, any) {},
			takesContext: true,
		},
		{
			name:         "2 arg, 1 returns",
			handler:      func(context.Context, any) error { return nil },
			takesContext: true,
		},
		{
			name:         "2 arg, 2 returns",
			handler:      func(context.Context, any) (any, error) { return 1, nil },
			takesContext: true,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			actual := "unexpected"
			logFatalf = func(format string, v ...interface{}) {
				actual = fmt.Sprintf(format, v...)
			}
			switch h := testCase.handler.(type) {
			case func():
				StartWithOptionsTypeSafe[any, any](h)
			case func() error:
				StartWithOptionsTypeSafe[any, any](h)
			case func(any):
				StartWithOptionsTypeSafe[any, any](h)
			case func(any) error:
				StartWithOptionsTypeSafe[any, any](h)
			case func() (any, error):
				StartWithOptionsTypeSafe[any, any](h)
			case func(any) (any, error):
				StartWithOptionsTypeSafe[any, any](h)
			case func(context.Context, any):
				StartWithOptionsTypeSafe[any, any](h)
			case func(context.Context, any) error:
				StartWithOptionsTypeSafe[any, any](h)
			case func(context.Context, any) (any, error):
				StartWithOptionsTypeSafe[any, any](h)
			default:
				assert.Fail(t, "Unexpected type: %T for test case: %s", h, testCase.name)
			}

			assert.Equal(t, "expected AWS Lambda environment variables [_LAMBDA_SERVER_PORT AWS_LAMBDA_RUNTIME_API] are not defined", actual)

			handlerType := reflect.TypeOf(testCase.handler)

			handlerTakesContext, err := validateArguments(handlerType)
			assert.NoError(t, err)
			assert.Equal(t, testCase.takesContext, handlerTakesContext)

			err = validateReturns(handlerType)
			assert.NoError(t, err)
		})
	}
}
