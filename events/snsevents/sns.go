// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package snsevents

import (
	"time"
)

type SnsEvent struct {
	Records []SnsEventRecord `json:"Records"`
}

type SnsEventRecord struct {
	EventVersion         string    `json:"EventVersion"`
	EventSubscriptionArn string    `json:"EventSubscriptionArn"`
	EventSource          string    `json:"EventSource"`
	Sns                  SnsEntity `json:"Sns"`
}

type SnsEntity struct {
	Signature         string                 `json:"Signature"`
	MessageId         string                 `json:"MessageId"`
	Type              string                 `json:"Type"`
	TopicArn          string                 `json:"TopicArn"`
	MessageAttributes map[string]interface{} `json:"MessageAttributes"`
	SignatureVersion  string                 `json:"SignatureVersion"`
	Timestamp         time.Time              `json:"Timestamp"`
	SigningCertUrl    string                 `json:"SigningCertUrl"`
	Message           string                 `json:"Message"`
	UnsubscribeUrl    string                 `json:"UnsubscribeUrl"`
	Subject           string                 `json:"Subject"`
}
