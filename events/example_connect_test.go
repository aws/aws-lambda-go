package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleConnectEvent() {
	lambda.Start(func(ctx context.Context, connectEvent *events.ConnectEvent) (*events.ConnectResponse, error) {
		fmt.Printf("Processing Connect event with ContactID %s.\n", connectEvent.Details.ContactData.ContactID)

		fmt.Printf("Invoked with %d parameters\n", len(connectEvent.Details.Parameters))
		for k, v := range connectEvent.Details.Parameters {
			fmt.Printf("%s : %s\n", k, v)
		}

		resp := &events.ConnectResponse{
			"Result":       "Success",
			"NewAttribute": "NewValue",
		}

		return resp, nil
	})
}
