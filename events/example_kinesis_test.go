package events_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleKinesisEvent() {
	lambda.Start(func(ctx context.Context, kinesisEvent *events.KinesisEvent) error {
		for _, record := range kinesisEvent.Records {
			kinesisRecord := record.Kinesis
			dataBytes := kinesisRecord.Data
			dataText := string(dataBytes)

			fmt.Printf("%s Data = %s \n", record.EventName, dataText)
		}

		return nil
	})
}

// This example transforms Kinesis Firehose records by converting the data to uppercase.
func ExampleKinesisFirehoseEvent() {
	lambda.Start(func(evnt *events.KinesisFirehoseEvent) (*events.KinesisFirehoseResponse, error) {
		fmt.Printf("InvocationID: %s\n", evnt.InvocationID)
		fmt.Printf("DeliveryStreamArn: %s\n", evnt.DeliveryStreamArn)
		fmt.Printf("Region: %s\n", evnt.Region)

		response := &events.KinesisFirehoseResponse{}

		for _, record := range evnt.Records {
			fmt.Printf("RecordID: %s\n", record.RecordID)
			fmt.Printf("ApproximateArrivalTimestamp: %s\n", record.ApproximateArrivalTimestamp)

			var transformedRecord events.KinesisFirehoseResponseRecord
			transformedRecord.RecordID = record.RecordID
			transformedRecord.Result = events.KinesisFirehoseTransformedStateOk
			transformedRecord.Data = []byte(strings.ToUpper(string(record.Data)))

			response.Records = append(response.Records, transformedRecord)
		}

		return response, nil
	})
}

func ExampleKinesisAnalyticsOutputDeliveryEvent() {
	lambda.Start(func(ctx context.Context, kinesisAnalyticsEvent *events.KinesisAnalyticsOutputDeliveryEvent) (*events.KinesisAnalyticsOutputDeliveryResponse, error) {
		responses := &events.KinesisAnalyticsOutputDeliveryResponse{}
		responses.Records = make([]events.KinesisAnalyticsOutputDeliveryResponseRecord, len(kinesisAnalyticsEvent.Records))

		for i, record := range kinesisAnalyticsEvent.Records {
			responses.Records[i] = events.KinesisAnalyticsOutputDeliveryResponseRecord{
				RecordID: record.RecordID,
				Result:   events.KinesisAnalyticsOutputDeliveryOK,
			}

			dataBytes := record.Data
			dataText := string(dataBytes)

			fmt.Printf("%s Data = %s \n", record.RecordID, dataText)
		}
		return responses, nil
	})
}
