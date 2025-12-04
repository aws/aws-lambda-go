package events_test

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleLambdaFunctionURLStreamingResponse() {
	lambda.Start(func() (*events.LambdaFunctionURLStreamingResponse, error) {
		return &events.LambdaFunctionURLStreamingResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "text/html",
			},
			Body: strings.NewReader("<html><body>Hello World!</body></html>"),
		}, nil
	})
}
