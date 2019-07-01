// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

type S3BatchJobEvent struct {
	InvocationSchemaVersion string           `json:"invocationSchemaVersion"`
	InvocationID            string           `json:"invocationId"`
	Job                     S3BatchJob       `json:"job"`
	Tasks                   []S3BatchJobTask `json:"tasks"`
}

type S3BatchJob struct {
	ID string `json:"id"`
}

type S3BatchJobTask struct {
	TaskID      string `json:"taskId"`
	S3Key       string `json:"s3Key"`
	S3VersionID string `json:"s3VersionId"`
	S3BucketARN string `json:"s3BucketArn"`
}

type S3BatchJobResponse struct {
	InvocationSchemaVersion string             `json:"invocationSchemaVersion"`
	TreatMissingKeysAs      string             `json:"treatMissingKeysAs"`
	InvocationID            string             `json:"invocationId"`
	Results                 []S3BatchJobResult `json:"results"`
}

type S3BatchJobResult struct {
	TaskID       string `json:"taskId"`
	ResultCode   string `json:"resultCode"`
	ResultString string `json:"resultString"`
}
