// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package kinesisfhevents

import "github.com/aws/aws-lambda-go/events"

// KinesisFirehoseEvent represents the input event from Amazon Kinesis Firehose. It is used as the input parameter.
type KinesisFirehoseEvent struct {
	InvocationId      string                `json:"invocationId"`
	DeliveryStreamArn string                `json:"deliveryStreamArn"`
	Region            string                `json:"region"`
	Records           []FirehoseEventRecord `json:"records"`
}

type FirehoseEventRecord struct {
	RecordId                    string                       `json:"recordId"`
	ApproximateArrivalTimestamp events.MilliSecondsEpochTime `json:"approximateArrivalTimestamp"`
	Data                        []byte                       `json:"data"`
}

// Constants used for describing the transformation result
const (
	TransformedStateOk               = "TRANSFORMED_STATE_OK"
	TransformedStateDropped          = "TRANSFORMED_STATE_DROPPED"
	TransformedStateProcessingFailed = "TRANSFORMED_STATE_PROCESSINGFAILED"
)

type KinesisFirehoseResponse struct {
	Records []FirehoseResponseRecord `json:"records"`
}

type FirehoseResponseRecord struct {
	RecordId string `json:"recordId"`
	Result   string `json:"result"` // The status of the transformation. May be TransformedStateOk, TransformedStateDropped or TransformedStateProcessingFailed
	Data     []byte `json:"data"`
}
