package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleSQSEvent() {
	lambda.Start(func(ctx context.Context, sqsEvent *events.SQSEvent) error {
		for _, message := range sqsEvent.Records {
			fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		}
		return nil
	})
}
