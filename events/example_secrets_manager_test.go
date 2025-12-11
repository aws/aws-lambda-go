package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// This Lambda function handles a SecretsManager secret rotation event.
// The rotation process has four steps: createSecret, setSecret, testSecret, and finishSecret.
func ExampleSecretsManagerSecretRotationEvent() {
	lambda.Start(func(ctx context.Context, event *events.SecretsManagerSecretRotationEvent) error {
		fmt.Printf("rotating secret %s with token %s\n",
			event.SecretID, event.ClientRequestToken)

		switch event.Step {
		case "createSecret":
			// create
		case "setSecret":
			// set
		case "finishSecret":
			// finish
		case "testSecret":
			// test
		}

		return nil
	})
}
