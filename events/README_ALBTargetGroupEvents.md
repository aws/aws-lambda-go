# Overview

Lambda Target Group events consist of a request that was routed to a Lambda function by a Lambda Target Group. When this happens, the Target Group expects the result of the function to be the response that the Target Group should respond with.

# Sample Function

The following is a sample class and Lambda function that receives AWS Lambda Target Group event as an input, writes some of the incoming data to CloudWatch Logs, and responds with a 200 status and the same body as the request. (Note that by default anything written to Console will be logged as CloudWatch Logs events.)

```go

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	fmt.Printf("Processing request data for traceId %s.\n", request.Headers["x-amzn-trace-id"])
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	return events.ALBTargetGroupResponse{Body: request.Body, StatusCode: 200, StatusDescription: "200 OK", IsBase64Encoded: false}, nil
}

func main() {
	lambda.Start(handleRequest)
}
```
