package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Stream notifications are delivered to the Lambda handler whenever data in the DynamoDB table is modified.
// Depending on the Stream settings, a StreamRecord may contain:
//   - Keys: key attributes of the modified item
//   - NewImage: the entire item, as it appears after it was modified
//   - OldImage: the entire item, as it appeared before it was modified
func ExampleDynamoDBEvent() {
	lambda.Start(func(ctx context.Context, e *events.DynamoDBEvent) {
		for _, record := range e.Records {
			fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

			for name, value := range record.Change.NewImage {
				if value.DataType() == events.DataTypeString {
					fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
				}
			}
		}
	})
}
