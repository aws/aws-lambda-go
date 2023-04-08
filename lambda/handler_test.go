// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/lambda/handlertrace"
	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidHandlers(t *testing.T) {
	type valuer interface {
		Value(key interface{}) interface{}
	}

	type customContext interface {
		context.Context
		MyCustomMethod()
	}

	type myContext interface {
		Deadline() (deadline time.Time, ok bool)
		Done() <-chan struct{}
		Err() error
		Value(key interface{}) interface{}
	}

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
			//nolint: stylecheck
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
		{
			name:     "the handler takes the empty interface",
			expected: nil,
			handler: func(v interface{}) error {
				if _, ok := v.(context.Context); ok {
					return errors.New("v should not be a Context")
				}
				return nil
			},
		},
		{
			name:     "the handler takes a subset of context.Context",
			expected: errors.New("handler takes an interface, but it is not context.Context: \"valuer\""),
			handler: func(ctx valuer) {
			},
		},
		{
			name:     "the handler takes a same interface with context.Context",
			expected: nil,
			handler: func(ctx myContext) {
			},
		},
		{
			name:     "the handler takes a superset of context.Context",
			expected: errors.New("handler takes an interface, but it is not context.Context: \"customContext\""),
			handler: func(ctx customContext) {
			},
		},
		{
			name:     "the handler takes two arguments and first argument is a subset of context.Context",
			expected: errors.New("handler takes two arguments, but the first is not Context. got interface"),
			handler: func(ctx valuer, v interface{}) {
			},
		},
		{
			name:     "the handler takes two arguments and first argument is a same interface with context.Context",
			expected: nil,
			handler: func(ctx myContext, v interface{}) {
			},
		},
		{
			name:     "the handler takes two arguments and first argument is a superset of context.Context",
			expected: errors.New("handler takes two arguments, but the first is not Context. got interface"),
			handler: func(ctx customContext, v interface{}) {
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			lambdaHandler := NewHandler(testCase.handler)
			_, err := lambdaHandler.Invoke(context.TODO(), []byte("{}"))
			assert.Equal(t, testCase.expected, err)
		})
	}
}

type arbitraryJSON struct {
	json []byte
	err  error
}

func (a arbitraryJSON) MarshalJSON() ([]byte, error) {
	return a.json, a.err
}

type staticHandler struct {
	body []byte
}

func (h *staticHandler) Invoke(_ context.Context, _ []byte) ([]byte, error) {
	return h.body, nil
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
		options  []Option
	}{
		{
			input:    `"Lambda"`,
			expected: expected{`null`, nil},
			handler:  func(_ string) {},
		},
		{
			input:    `"Lambda"`,
			expected: expected{`"Hello Lambda!"`, nil},
			handler: func(name string) (string, error) {
				return hello(name), nil
			},
		},
		{
			expected: expected{`"Hello No Value!"`, nil},
			handler: func(ctx context.Context) (string, error) {
				return hello("No Value"), nil
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
		{
			name:     "WithSetEscapeHTML(false)",
			expected: expected{`"<html><body>html in json string!</body></html>"`, nil},
			handler: func() (string, error) {
				return "<html><body>html in json string!</body></html>", nil
			},
			options: []Option{WithSetEscapeHTML(false)},
		},
		{
			name:     "WithSetEscapeHTML(true)",
			expected: expected{`"\u003chtml\u003e\u003cbody\u003ehtml in json string!\u003c/body\u003e\u003c/html\u003e"`, nil},
			handler: func() (string, error) {
				return "<html><body>html in json string!</body></html>", nil
			},
			options: []Option{WithSetEscapeHTML(true)},
		},
		{
			name:     `WithSetIndent(">>", "  ")`,
			expected: expected{"{\n>>  \"Foo\": \"Bar\"\n>>}\n", nil},
			handler: func() (interface{}, error) {
				return struct{ Foo string }{"Bar"}, nil
			},
			options: []Option{WithSetIndent(">>", "  ")},
		},
		{
			name:     "bytes are base64 encoded strings",
			input:    `"aGVsbG8="`,
			expected: expected{`"aGVsbG95b2xv"`, nil},
			handler: func(_ context.Context, req []byte) ([]byte, error) {
				return append(req, []byte("yolo")...), nil
			},
		},
		{
			name:     "Handler interface implementations are passthrough",
			expected: expected{`<xml>hello</xml>`, nil},
			handler:  &staticHandler{body: []byte(`<xml>hello</xml>`)},
		},
		{
			name:     "io.Reader responses are passthrough",
			expected: expected{`<yolo>yolo</yolo>`, nil},
			handler: func() (io.Reader, error) {
				return strings.NewReader(`<yolo>yolo</yolo>`), nil
			},
		},
		{
			name:     "io.Reader responses that are byte buffers are passthrough",
			expected: expected{`<yolo>yolo</yolo>`, nil},
			handler: func() (*bytes.Buffer, error) {
				return bytes.NewBuffer([]byte(`<yolo>yolo</yolo>`)), nil
			},
		},
		{
			name:     "io.Reader responses that are also json serializable, handler returns the json, ignoring the reader",
			expected: expected{`{"Yolo":"yolo"}`, nil},
			handler: func() (io.Reader, error) {
				return struct {
					io.Reader `json:"-"`
					Yolo      string
				}{
					Reader: strings.NewReader(`<yolo>yolo</yolo>`),
					Yolo:   "yolo",
				}, nil
			},
		},
		{
			name:     "types that are not json serializable result in an error",
			expected: expected{``, errors.New("json: error calling MarshalJSON for type struct { lambda.arbitraryJSON }: barf")},
			handler: func() (interface{}, error) {
				return struct {
					arbitraryJSON
				}{
					arbitraryJSON{nil, errors.New("barf")},
				}, nil
			},
		},
		{
			name:     "io.Reader responses that not json serializable remain passthrough",
			expected: expected{`wat`, nil},
			handler: func() (io.Reader, error) {
				return struct {
					arbitraryJSON
					io.Reader
				}{
					arbitraryJSON{nil, errors.New("barf")},
					strings.NewReader("wat"),
				}, nil
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			lambdaHandler := newHandler(testCase.handler, testCase.options...)
			t.Run("via Handler.Invoke", func(t *testing.T) {
				response, err := lambdaHandler.Invoke(context.TODO(), []byte(testCase.input))
				if testCase.expected.err != nil {
					assert.EqualError(t, err, testCase.expected.err.Error())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, testCase.expected.val, string(response))
				}
			})
			t.Run("via handlerOptions.handlerFunc", func(t *testing.T) {
				response, err := lambdaHandler.handlerFunc(context.TODO(), []byte(testCase.input))
				if testCase.expected.err != nil {
					assert.EqualError(t, err, testCase.expected.err.Error())
				} else {
					assert.NoError(t, err)
					require.NotNil(t, response)
					responseBytes, err := io.ReadAll(response)
					assert.NoError(t, err)
					assert.Equal(t, testCase.expected.val, string(responseBytes))
				}
			})

		})
	}
}

func TestInvalidJsonInput(t *testing.T) {
	lambdaHandler := NewHandler(func(s string) error { return nil })
	_, err := lambdaHandler.Invoke(context.TODO(), []byte(`{"invalid json`))
	assert.Equal(t, "unexpected EOF", err.Error())
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
