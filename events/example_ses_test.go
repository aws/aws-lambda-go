package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleSimpleEmailEvent() {
	lambda.Start(func(ctx context.Context, sesEvent *events.SimpleEmailEvent) error {
		for _, record := range sesEvent.Records {
			ses := record.SES
			fmt.Printf("[%s - %s] Mail = %+v, Receipt = %+v \n", record.EventVersion, record.EventSource, ses.Mail, ses.Receipt)
		}

		return nil
	})
}
