package events_test

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleAPIGatewayProxyStreamingResponse() {
	lambda.Start(func() (*events.APIGatewayProxyStreamingResponse, error) {
		return &events.APIGatewayProxyStreamingResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "text/html",
			},
			Body: strings.NewReader("<html><body>Hello World!</body></html>"),
		}, nil
	})
}
