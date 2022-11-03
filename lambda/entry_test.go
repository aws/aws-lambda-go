// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ctxTestKey struct{}

func TestStartRuntimeAPIWithContext(t *testing.T) {
	server, _ := runtimeAPIServer("null", 1) // serve a single invoke, and then cause an internal error
	expected := "expected"
	actual := "unexpected"

	os.Setenv("AWS_LAMBDA_RUNTIME_API", strings.Split(server.URL, "://")[1])
	defer os.Unsetenv("AWS_LAMBDA_RUNTIME_API")
	logFatalf = func(format string, v ...interface{}) {}
	defer func() { logFatalf = log.Fatalf }()

	StartWithContext(context.WithValue(context.Background(), ctxTestKey{}, expected), func(ctx context.Context) error {
		actual, _ = ctx.Value(ctxTestKey{}).(string)
		return nil
	})

	assert.Equal(t, expected, actual)
}

func TestStartWithOptionsTypeSafe(t *testing.T) {
	testCases := []struct {
		name    string
		handler any
	}{
		{
			name:    "0 arg, 0 returns",
			handler: func() {},
		},
		{
			name:    "0 arg, 1 returns",
			handler: func() error { return nil },
		},
		{
			name:    "1 arg, 0 returns",
			handler: func(any) {},
		},
		{
			name:    "1 arg, 1 returns",
			handler: func(any) error { return nil },
		},
		{
			name:    "0 arg, 2 returns",
			handler: func() (any, error) { return 1, nil },
		},
		{
			name:    "1 arg, 2 returns",
			handler: func(any) (any, error) { return 1, nil },
		},
		{
			name:    "2 arg, 0 returns",
			handler: func(context.Context, any) {},
		},
		{
			name:    "2 arg, 1 returns",
			handler: func(context.Context, any) error { return nil },
		},
		{
			name:    "2 arg, 2 returns",
			handler: func(context.Context, any) (any, error) { return 1, nil },
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

			_, err := validateArguments(handlerType)
			assert.NoError(t, err)

			err = validateReturns(handlerType)
			assert.NoError(t, err)
		})
	}
}
