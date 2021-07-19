// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/lambda/handlertrace"
	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/stretchr/testify/assert"
)

func TestInvalidHandlers(t *testing.T) {

	testCases := []struct {
		name     string
		handler  interface{}
		expected error
	}{
		{
			name:     "nil handler",
			expected: errors.New("handler is nil"),
			handler:  nil,
		},
		{
			name:     "handler is not a function",
			expected: errors.New("handler kind struct is not func"),
			handler:  struct{}{},
		},
		{
			name:     "handler declares too many arguments",
			expected: errors.New("handlers may not take more than two arguments, but handler takes 3"),
			handler: func(n context.Context, x string, y string) error {
				return nil
			},
		},
		{
			name:     "two argument handler does not context as first argument",
			expected: errors.New("handler takes two arguments, but the first is not Context. got string"),
			handler: func(a string, x context.Context) error {
				return nil
			},
		},
		{
			name:     "handler returns too many values",
			expected: errors.New("handler may not return more than two values"),
			handler: func() (error, error, error) {
				return nil, nil, nil
			},
		},
		{
			name:     "handler returning two values does not declare error as the second return value",
			expected: errors.New("handler returns two values, but the second does not implement error"),
			handler: func() (error, string) {
				return nil, "hello"
			},
		},
		{
			name:     "handler returning a single value does not implement error",
			expected: errors.New("handler returns a single value, but it does not implement error"),
			handler: func() string {
				return "hello"
			},
		},
		{
			name:     "no return value should not result in error",
			expected: nil,
			handler: func() {
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			lambdaHandler := NewHandler(testCase.handler)
			_, err := lambdaHandler.Invoke(context.TODO(), make([]byte, 0))
			assert.Equal(t, testCase.expected, err)
		})
	}
}

type expected struct {
	val string
	err error
}

func TestInvokes(t *testing.T) {
	hello := func(s string) string {
		return fmt.Sprintf("Hello %s!", s)
	}
	hellop := func(s *string) *string {
		v := hello(*s)
		return &v
	}

	testCases := []struct {
		name     string
		input    string
		expected expected
		handler  interface{}
	}{
		{
			input:    `"Lambda"`,
			expected: expected{`"Hello Lambda!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{`"Hello Lambda!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{`"Hello Lambda!"`, nil},
			handler: func(ctx context.Context, name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{`"Hello Lambda!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{`"Hello Lambda!"`, nil},
			handler: func(name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{`"Hello Lambda!"`, nil},
			handler: func(ctx context.Context, name *string) (*string, error) {
				return hellop(name), nil
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{"", errors.New("bad stuff")},
			handler: func() error {
				return errors.New("bad stuff")
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{"", errors.New("bad stuff")},
			handler: func() (interface{}, error) {
				return nil, errors.New("bad stuff")
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{"", errors.New("bad stuff")},
			handler: func(e interface{}) (interface{}, error) {
				return nil, errors.New("bad stuff")
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{"", errors.New("bad stuff")},
			handler: func(ctx context.Context, e interface{}) (interface{}, error) {
				return nil, errors.New("bad stuff")
			},
		},
		{
			name:     "basic input struct serialization",
			input:    `{"custom":9001}`,
			expected: expected{`9001`, nil},
			handler: func(event struct{ Custom int }) (int, error) {
				return event.Custom, nil
			},
		},
		{
			name:     "basic output struct serialization",
			input:    `9001`,
			expected: expected{`{"Number":9001}`, nil},
			handler: func(event int) (struct{ Number int }, error) {
				return struct{ Number int }{event}, nil
			},
		},
		{
			input:    `"Lambda"`,
			expected: expected{"", messages.InvokeResponse_Error{Message: "message", Type: "type"}},
			handler: func(e interface{}) (interface{}, error) {
				return nil, messages.InvokeResponse_Error{Message: "message", Type: "type"}
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			lambdaHandler := NewHandler(testCase.handler)
			response, err := lambdaHandler.Invoke(context.TODO(), []byte(testCase.input))
			if testCase.expected.err != nil {
				assert.Equal(t, testCase.expected.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected.val, string(response))
			}
		})
	}
}

func TestInvalidJsonInput(t *testing.T) {
	lambdaHandler := NewHandler(func(s string) error { return nil })
	_, err := lambdaHandler.Invoke(context.TODO(), []byte(`{"invalid json`))
	assert.Equal(t, "unexpected end of JSON input", err.Error())
}

func TestHandlerTrace(t *testing.T) {
	handler := NewHandler(func(ctx context.Context, x int) (int, error) {
		if x != 123 {
			t.Error(x)
		}
		return 456, nil
	})
	requestHistory := ""
	responseHistory := ""
	checkInt := func(e interface{}, expected int) {
		nt, ok := e.(int)
		if !ok {
			t.Error("not int as expected", e)
			return
		}
		if nt != expected {
			t.Error("unexpected value", nt, expected)
		}
	}
	ctx := context.Background()
	ctx = handlertrace.NewContext(ctx, handlertrace.HandlerTrace{}) // empty HandlerTrace
	ctx = handlertrace.NewContext(ctx, handlertrace.HandlerTrace{   // with RequestEvent
		RequestEvent: func(c context.Context, e interface{}) {
			requestHistory += "A"
			checkInt(e, 123)
		},
	})
	ctx = handlertrace.NewContext(ctx, handlertrace.HandlerTrace{ // with ResponseEvent
		ResponseEvent: func(c context.Context, e interface{}) {
			responseHistory += "X"
			checkInt(e, 456)
		},
	})
	ctx = handlertrace.NewContext(ctx, handlertrace.HandlerTrace{ // with RequestEvent and ResponseEvent
		RequestEvent: func(c context.Context, e interface{}) {
			requestHistory += "B"
			checkInt(e, 123)
		},
		ResponseEvent: func(c context.Context, e interface{}) {
			responseHistory += "Y"
			checkInt(e, 456)
		},
	})
	ctx = handlertrace.NewContext(ctx, handlertrace.HandlerTrace{}) // empty HandlerTrace

	payload := []byte(`123`)
	js, err := handler.Invoke(ctx, payload)
	if err != nil {
		t.Error("unexpected handler error", err)
	}
	if string(js) != "456" {
		t.Error("unexpected handler output", string(js))
	}
	if requestHistory != "AB" {
		t.Error("request callbacks not called as expected", requestHistory)
	}
	if responseHistory != "XY" {
		t.Error("response callbacks not called as expected", responseHistory)
	}
}
