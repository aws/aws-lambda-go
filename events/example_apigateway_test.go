package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// API Gateway events consist of a request that was routed to a Lambda function by API Gateway.
// When this happens, API Gateway expects the result of the function to be the response that API Gateway should respond with.
func ExampleAPIGatewayProxyRequest() {
	lambda.Start(func(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
		fmt.Printf("Body size = %d.\n", len(request.Body))

		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("    %s: %s\n", key, value)
		}

		return &events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
	})
}
