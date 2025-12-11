package cfn_test

import (
	"context"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
)

// CloudFormation custom resources require a different response handling due to the way stacks execute.
// The cfn.LambdaWrap helper catches all errors and ensures the correct response is sent to the
// pre-signed URL that comes with the event.
//
// This example safely 'Echo' back anything given into the Echo parameter within the Custom Resource call.
func Example() {
	lambda.Start(cfn.LambdaWrap(func(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
		v, _ := event.ResourceProperties["Echo"].(string)

		data = map[string]interface{}{
			"Echo": v,
		}

		return
	}))
}
