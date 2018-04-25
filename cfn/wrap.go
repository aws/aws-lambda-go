// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package cfn

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

// CustomResourceLambdaFunction is a standard form Lambda for a Custom Resource.
type CustomResourceLambdaFunction func(context.Context, Event) (reason string, err error)

// CustomResourceFunction is a representation of the customer's Custom Resource function.
// LambdaWrap will take the returned values and turn them into a response to be sent
// to CloudFormation.
type CustomResourceFunction func(context.Context, Event) (physicalResourceID string, data map[string]interface{}, err error)

func lambdaWrapWithClient(lambdaFunction CustomResourceFunction, client httpClient) (fn CustomResourceLambdaFunction) {
	fn = func(ctx context.Context, event Event) (reason string, err error) {
		r := NewResponse(&event)
		r.PhysicalResourceID, r.Data, err = lambdaFunction(ctx, event)

		if err != nil {
			r.Status = StatusFailed
			r.Reason = err.Error()
			log.Printf("sending status failed: %s", r.Reason)
		} else {
			r.Status = StatusSuccess

			if r.PhysicalResourceID == "" {
				log.Println("PhysicalResourceID must exist on creation, copying Log Stream name")
				r.PhysicalResourceID = lambdacontext.LogStreamName
			}
		}

		err = r.sendWith(client)
		if err != nil {
			reason = err.Error()
		}

		return
	}

	return
}

// LambdaWrap returns a CustomResourceLambdaFunction which is something lambda.Start()
// will understand. The purpose of doing this is so that Response Handling boiler
// plate is taken away from the customer and it makes writing a Custom Resource
// simpler.
//
// 	func myLambda(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
// 		physicalResourceID = "arn:...."
// 		return
// 	}
//
// 	func main() {
// 		lambda.Start(cfn.LambdaWrap(myLambda))
// 	}
func LambdaWrap(lambdaFunction CustomResourceFunction) (fn CustomResourceLambdaFunction) {
	return lambdaWrapWithClient(lambdaFunction, http.DefaultClient)
}
