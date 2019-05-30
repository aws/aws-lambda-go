// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"time"
)

type SNSEvent struct {
	Records []SNSEventRecord `json:"Records"`
}

type SNSEventRecord struct {
	EventVersion         string    `json:"EventVersion"`
	EventSubscriptionArn string    `json:"EventSubscriptionArn"`
	EventSource          string    `json:"EventSource"`
	SNS                  SNSEntity `json:"Sns"`
}

type SNSEntity struct {
	Signature         string                 `json:"Signature"`
	MessageID         string                 `json:"MessageId"`
	Type              string                 `json:"Type"`
	TopicArn          string                 `json:"TopicArn"`
	MessageAttributes map[string]interface{} `json:"MessageAttributes"`
	SignatureVersion  string                 `json:"SignatureVersion"`
	Timestamp         time.Time              `json:"Timestamp"`
	SigningCertURL    string                 `json:"SigningCertUrl"`
	Message           string                 `json:"Message"`
	UnsubscribeURL    string                 `json:"UnsubscribeUrl"`
	Subject           string                 `json:"Subject"`
}

type CloudFormationAlarm struct {
	AlarmName        string `json:"AlarmName"`
	AlarmDescription string `json:"AlarmDescription"`
	AWSAccountId     string `json:"AWSAccountId"`

	NewStateValue   string `json:"NewStateValue"`
	NewStateReason  string `json:"NewStateReason"`
	OldStateValue   string `json:"OldStateValue"`
	StateChangeTime string `json:"StateChangeTime"`
	Region          string `json:"Region"`
}

type CloudFormationTrigger struct {
	MetricName                       string                           `json:"MetricName"`
	Namespace                        string                           `json:"Namespace"`
	StatisticType                    string                           `json:"StatisticType"`
	Statistic                        string                           `json:"Statistic"`
	Unit                             string                           `json:"Unit"`
	Dimensions                       []CloudFormationTriggerDimension `json:"Dimensions"`
	Period                           int64                            `json:"Period"`
	EvaluationPeriods                int64                            `json:"EvaluationPeriods"`
	ComparisonOperator               string                           `json:"ComparisonOperator"`
	Threshold                        int64                            `json:"Threshold"`
	TreatMissingData                 string                           `json:"TreatMissingData"`
	EvaluateLowSampleCountPercentile string                           `json:"EvaluateLowSampleCountPercentile"`
}

type CloudFormationTriggerDimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
