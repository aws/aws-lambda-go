# Overview

This package provides input types for Lambda functions that process Amazon API Gateway events.

API Gateway events consist of a request that was routed to a Lambda function by API Gateway. When this happens, API Gateway expects the result of the function to be the response that API Gateway should respond with.

# Sample Function

The following is a sample class and Lambda function that receives Amazon API Gateway event record data as an input, writes some of the record data to CloudWatch Logs, and responds with a 200 status and the same body as the request. (Note that by default anything written to Console will be logged as CloudWatch Logs events.)

```go

import (
    "strings"
    "github.com/aws/aws-lambda-go/events/apigwevents"
)

func handleRequest(ctx context.Context, request apigwevents.ApiGatewayProxyRequest) apigwevents.ApiGatewayProxyResponse {
    fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestId)
    fmt.Printf("Body size = %d.\n", len(request.Body))

    fmt.Println("Headers:")
    for key, value := range request.Headers {
        fmt.Printf("    %s: %s\n", key, value)
    }

    return ApiGatewayProxyResponse { Body: request.Body, StatusCode: 200 }
}
```
