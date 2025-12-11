package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ALB Target Group events consist of a request that was routed to a Lambda function
// which is a registered target of an Application Load Balancer Target Group.
// When this happens, ALB expects the result of the function to be the response that ALB should respond with.
func ExampleALBTargetGroupRequest() {
	lambda.Start(func(ctx context.Context, request *events.ALBTargetGroupRequest) (*events.ALBTargetGroupResponse, error) {
		fmt.Printf("Processing request data for traceId %s.\n", request.Headers["x-amzn-trace-id"])
		fmt.Printf("Body size = %d.\n", len(request.Body))

		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("    %s: %s\n", key, value)
		}

		return &events.ALBTargetGroupResponse{Body: request.Body, StatusCode: 200, StatusDescription: "200 OK", IsBase64Encoded: false, Headers: map[string]string{}}, nil
	})
}
