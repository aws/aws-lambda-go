// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main () {
	lambda.Start(func (ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
		return &events.APIGatewayProxyResponse{
			Body: fmt.Sprintf("Hello %v", event),
		}, nil
	})
}
