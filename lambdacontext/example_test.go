package lambdacontext_test

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func ExampleFromContext() {
	lambda.Start(func(ctx context.Context) (string, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		log.Printf("Request ID: %s", lc.AwsRequestID)
		log.Printf("Function ARN: %s", lc.InvokedFunctionArn)
		return "success", nil
	})
}
