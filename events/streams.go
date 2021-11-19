package events

// StreamsEventResponse is the outer structure to report batch item failures for KinesisEvent and DynamoDBEvent.
type StreamsEventResponse struct {
	BatchItemFailures []BatchItemFailure `json:"batchItemFailures"`
}

// Record which failed processing.
type BatchItemFailure struct {
	ItemIdentifier string `json:"itemIdentifier"`
}
