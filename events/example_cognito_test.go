package events_test

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleCognitoEvent() {
	lambda.Start(func(cognitoEvent *events.CognitoEvent) error {
		for datasetName, datasetRecord := range cognitoEvent.DatasetRecords {
			fmt.Printf("[%s -- %s] %s -> %s -> %s \n",
				cognitoEvent.EventType,
				datasetName,
				datasetRecord.OldValue,
				datasetRecord.Op,
				datasetRecord.NewValue)
		}
		return nil
	})
}

// For setting up Cognito User Pools triggers, see:
// https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html
func ExampleCognitoEventUserPoolsPreSignup() {
	lambda.Start(func(event *events.CognitoEventUserPoolsPreSignup) (*events.CognitoEventUserPoolsPreSignup, error) {
		fmt.Printf("PreSignup of user: %s\n", event.UserName)
		event.Response.AutoConfirmUser = true
		return event, nil
	})
}

func ExampleCognitoEventUserPoolsPostConfirmation() {
	lambda.Start(func(event *events.CognitoEventUserPoolsPostConfirmation) (*events.CognitoEventUserPoolsPostConfirmation, error) {
		fmt.Printf("PostConfirmation for user: %s\n", event.UserName)
		return event, nil
	})
}

func ExampleCognitoEventUserPoolsPreAuthentication() {
	lambda.Start(func(event *events.CognitoEventUserPoolsPreAuthentication) (*events.CognitoEventUserPoolsPreAuthentication, error) {
		fmt.Printf("PreAuthentication of user: %s\n", event.UserName)
		return event, nil
	})
}

func ExampleCognitoEventUserPoolsPreTokenGen() {
	lambda.Start(func(event *events.CognitoEventUserPoolsPreTokenGen) (*events.CognitoEventUserPoolsPreTokenGen, error) {
		fmt.Printf("PreTokenGen of user: %s\n", event.UserName)
		event.Response.ClaimsOverrideDetails.ClaimsToSuppress = []string{"family_name"}
		return event, nil
	})
}

// These Lambda triggers issue and verify their own challenges as part of a user pool custom authentication flow.
// For setting up custom authentication, see:
// https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-challenge.html
func ExampleCognitoEventUserPoolsDefineAuthChallenge() {
	lambda.Start(func(event *events.CognitoEventUserPoolsDefineAuthChallenge) (*events.CognitoEventUserPoolsDefineAuthChallenge, error) {
		fmt.Printf("Define Auth Challenge: %+v\n", event)
		return event, nil
	})
}
