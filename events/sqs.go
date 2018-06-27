// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

type SQSEvent struct {
	Records []SQSEventRecord `json:"Records"`
}

type SQSEventRecord struct {
	Body              string                 `json:"body"`
	ReceiptHandle     string                 `json:"receiptHandle"`
	MD5OfBody         string                 `json:"md5OfBody"`
	EventSourceARN    string                 `json:"eventSourceARN"`
	EventSource       string                 `json:"eventSource"`
	AWSRegion         string                 `json:"awsRegion"`
	MessageID         string                 `json:"messageId"`
	Attributes        map[string]interface{} `json:"attributes"`
	MessageAttributes map[string]interface{} `json:"messageAttributes"`
}
