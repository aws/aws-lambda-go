package lambda_test

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func Example() {
	lambda.Start(func() (string, error) {
		return "Hello λ!", nil
	})
}

// Handlers can return io.Reader to stream response data.
// This example uses a pipe to send data in chunks with delays.
//
// See https://docs.aws.amazon.com/lambda/latest/dg/configuration-response-streaming.html
func Example_ioReader() {
	lambda.Start(func() (io.Reader, error) {
		r, w := io.Pipe()
		go func() {
			defer w.Close()
			_, _ = w.Write([]byte("<html><body>"))
			time.Sleep(100 * time.Millisecond)
			_, _ = w.Write([]byte("<h1>Hello</h1>"))
			time.Sleep(100 * time.Millisecond)
			_, _ = w.Write([]byte("<p>World!</p>"))
			time.Sleep(100 * time.Millisecond)
			_, _ = w.Write([]byte("</body></html>"))
		}()
		return r, nil
	})
}

func ExampleWithContext() {
	lambda.StartWithOptions(
		func(ctx context.Context) (string, error) {
			return ctx.Value("foo").(string), nil
		},
		lambda.WithContext(context.WithValue(context.Background(), "foo", "bar")),
	)
}

func ExampleWithContextValue() {
	lambda.StartWithOptions(
		func(ctx context.Context) (string, error) {
			return ctx.Value("foo").(string), nil
		},
		lambda.WithContextValue("foo", "bar"),
	)
}

func ExampleWithSetEscapeHTML() {
	lambda.StartWithOptions(
		func() (string, error) {
			return "<html><body>hello!</body></html>", nil
		},
		lambda.WithSetEscapeHTML(true),
	)
}

func ExampleWithSetIndent() {
	lambda.StartWithOptions(
		func(event interface{}) (interface{}, error) {
			return event, nil
		},
		lambda.WithSetIndent(">", " "),
	)
}

func ExampleWithUseNumber() {
	lambda.StartWithOptions(
		func(event interface{}) (interface{}, error) {
			return event, nil
		},
		lambda.WithUseNumber(true),
	)
}

func ExampleWithDisallowUnknownFields() {
	lambda.StartWithOptions(
		func(event interface{}) (interface{}, error) {
			return event, nil
		},
		lambda.WithDisallowUnknownFields(true),
	)
}

func ExampleWithEnableSIGTERM() {
	lambda.StartWithOptions(
		func(event interface{}) (interface{}, error) {
			return event, nil
		},
		lambda.WithEnableSIGTERM(func() {
			log.Print("function container shutting down...")
		}),
	)
}
