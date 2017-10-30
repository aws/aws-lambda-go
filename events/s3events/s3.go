// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package s3events

import (
	"time"
)

type S3Event struct {
	Records []S3EventRecord `json:"Records"`
}

type S3EventRecord struct {
	EventVersion      string              `json:"eventVersion"`
	EventSource       string              `json:"eventSource"`
	AwsRegion         string              `json:"awsRegion"`
	EventTime         time.Time           `json:"eventTime"`
	EventName         string              `json:"eventName"`
	PrincipalId       UserIdentity        `json:"userIdentity"`
	RequestParameters S3RequestParameters `json:"requestParameters"`
	ResponseElements  map[string]string   `json:"responseElements"`
	S3                S3Entity            `json:"s3"`
}

type UserIdentity struct {
	PrincipalId string `json:"principalId"`
}

type S3RequestParameters struct {
	SourceIpAddress string `json:"sourceIPAddress"`
}

type S3Entity struct {
	SchemaVersion   string   `json:"s3SchemaVersion"`
	ConfigurationId string   `json:"configurationId"`
	Bucket          S3Bucket `json:"bucket"`
	Object          S3Object `json:"object"`
}

type S3Bucket struct {
	Name          string       `json:"name"`
	OwnerIdentity UserIdentity `json:"ownerIdentity"`
	Arn           string       `json:"arn"`
}

type S3Object struct {
	Key           string `json:"key"`
	Size          int64  `json:"size"`
	UrlDecodedKey string `json:"urlDecodedKey"`
	VersionId     string `json:"versionId"`
	ETag          string `json:"eTag"`
	Sequencer     string `json:"sequencer"`
}
