// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package ddbevents

import "github.com/aws/aws-lambda-go/events"

// The DynamoDB stream event handled to Lambda
// http://docs.aws.amazon.com/lambda/latest/dg/eventsources.html#eventsources-ddb-update
type DynamoDbEvent struct {
	Records []DynamoDbEventRecord `json:"Records"`
}

type DynamoDbEventRecord struct {
	// The region in which the GetRecords request was received.
	AwsRegion string `json:"awsRegion"`

	// The main body of the stream record, containing all of the DynamoDB-specific
	// fields.
	Dynamodb StreamRecord `json:"dynamodb"`

	// A globally unique identifier for the event that was recorded in this stream
	// record.
	EventId string `json:"eventID"`

	// The type of data modification that was performed on the DynamoDB table:
	//
	//    * INSERT - a new item was added to the table.
	//
	//    * MODIFY - one or more of an existing item's attributes were modified.
	//
	//    * REMOVE - the item was deleted from the table
	EventName string `json:"eventName"`

	// The AWS service from which the stream record originated. For DynamoDB Streams,
	// this is aws:dynamodb.
	EventSource string `json:"eventSource"`

	// The version number of the stream record format. This number is updated whenever
	// the structure of Record is modified.
	//
	// Client applications must not assume that eventVersion will remain at a particular
	// value, as this number is subject to change at any time. In general, eventVersion
	// will only increase as the low-level DynamoDB Streams API evolves.
	EventVersion string `json:"eventVersion"`

	// The event source ARN of DynamoDB
	EventSourceArn string `json:"eventSourceARN"`
}

// A description of a single data modification that was performed on an item
// in a DynamoDB table.
type StreamRecord struct {

	// The approximate date and time when the stream record was created, in UNIX
	// epoch time (http://www.epochconverter.com/) format.
	ApproximateCreationDateTime events.SecondsEpochTime `json:"ApproximateCreationDateTime,omitempty"`

	// The primary key attribute(s) for the DynamoDB item that was modified.
	Keys map[string]AttributeValue `json:"Keys,omitempty"`

	// The item in the DynamoDB table as it appeared after it was modified.
	NewImage map[string]AttributeValue `json:"NewImage,omitempty"`

	// The item in the DynamoDB table as it appeared before it was modified.
	OldImage map[string]AttributeValue `json:"OldImage,omitempty"`

	// The sequence number of the stream record.
	SequenceNumber string `json:"SequenceNumber"`

	// The size of the stream record, in bytes.
	SizeBytes int64 `json:"SizeBytes"`

	// The type of data from the modified DynamoDB item that was captured in this
	// stream record.
	StreamViewType string `json:"StreamViewType"`
}

// Contains details about the type of identity that made the request.
type Identity struct {
	// A unique identifier for the entity that made the call. For Time To Live,
	// the principalId is "dynamodb.amazonaws.com".
	PrincipalId string `json:"PrincipalId"`

	// The type of the identity. For Time To Live, the type is "Service".
	Type string `json:"Type"`
}

const (
	KeyTypeHash  = "HASH"
	KeyTypeRange = "RANGE"
)

const (
	OperationTypeInsert = "INSERT"
	OperationTypeModify = "MODIFY"
	OperationTypeRemove = "REMOVE"
)

const (
	ShardIteratorTypeTrimHorizon         = "TRIM_HORIZON"
	ShardIteratorTypeLatest              = "LATEST"
	ShardIteratorTypeAtSequenceNumber    = "AT_SEQUENCE_NUMBER"
	ShardIteratorTypeAfterSequenceNumber = "AFTER_SEQUENCE_NUMBER"
)

const (
	StreamStatusEnabling  = "ENABLING"
	StreamStatusEnabled   = "ENABLED"
	StreamStatusDisabling = "DISABLING"
	StreamStatusDisabled  = "DISABLED"
)

const (
	StreamViewTypeNewImage        = "NEW_IMAGE"          // the entire item, as it appeared after it was modified.
	StreamViewTypeOldImage        = "OLD_IMAGE"          // the entire item, as it appeared before it was modified.
	StreamViewTypeNewAndOldImages = "NEW_AND_OLD_IMAGES" // both the new and the old item images of the item.
	StreamViewTypeKeysOnly        = "KEYS_ONLY"          // only the key attributes of the modified item.
)
