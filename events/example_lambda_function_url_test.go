package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleLambdaFunctionURLRequest() {
	lambda.Start(func(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
		fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
		fmt.Printf("Body size = %d.\n", len(request.Body))

		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("    %s: %s\n", key, value)
		}

		return &events.LambdaFunctionURLResponse{Body: request.Body, StatusCode: 200}, nil
	})
}
