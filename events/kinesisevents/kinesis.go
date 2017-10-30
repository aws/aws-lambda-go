// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package kinesisevents

import (
	"github.com/aws/aws-lambda-go/events"
)

type KinesisEvent struct {
	Records []KinesisEventRecord `json:"Records"`
}

type KinesisEventRecord struct {
	AwsRegion         string `json:"awsRegion"`
	EventId           string `json:"eventID"`
	EventName         string `json:"eventName"`
	EventSource       string `json:"eventSource"`
	EventSourceArn    string `json:"eventSourceARN"`
	EventVersion      string `json:"eventVersion"`
	InvokeIdentityArn string `json:"invokeIdentityArn"`
	Kinesis           Record `json:"kinesis"`
}

type Record struct {
	ApproximateArrivalTimestamp events.SecondsEpochTime `json:"approximateArrivalTimestamp"`
	Data                        []byte                  `json:"data"`
	EncryptionType              string                  `json:"encryptionType,omitempty"`
	PartitionKey                string                  `json:"partitionKey"`
	SequenceNumber              string                  `json:"sequenceNumber"`
	KinesisSchemaVersion        string                  `json:"kinesisSchemaVersion"`
}
