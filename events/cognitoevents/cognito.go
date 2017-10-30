// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package cognitoevents

// CognitoEvent contains data from an event sent from AWS Cognito
type CognitoEvent struct {
	DatasetName    string                   `json:"datasetName"`
	DatasetRecords map[string]DatasetRecord `json:"datasetRecords"`
	EventType      string                   `json:"eventType"`
	IdentityId     string                   `json:"identityId"`
	IdentityPoolId string                   `json:"identityPoolId"`
	Region         string                   `json:"region"`
	Version        int                      `json:"version"`
}

type DatasetRecord struct {
	NewValue string `json:"newValue"`
	OldValue string `json:"oldValue"`
	Op       string `json:"op"`
}
