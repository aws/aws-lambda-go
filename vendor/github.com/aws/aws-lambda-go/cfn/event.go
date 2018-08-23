// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package cfn

// RequestType represents the types of requests that
// come from a CloudFormation stack being run
type RequestType string

const (
	RequestCreate RequestType = "Create"
	RequestUpdate RequestType = "Update"
	RequestDelete RequestType = "Delete"
)

// Event is a representation of a Custom Resource
// request
type Event struct {
	RequestType           RequestType            `json:"RequestType"`
	RequestID             string                 `json:"RequestId"`
	ResponseURL           string                 `json:"ResponseURL"`
	ResourceType          string                 `json:"ResourceType"`
	PhysicalResourceID    string                 `json:"PhysicalResourceId,omitempty"`
	LogicalResourceID     string                 `json:"LogicalResourceId"`
	StackID               string                 `json:"StackId"`
	ResourceProperties    map[string]interface{} `json:"ResourceProperties"`
	OldResourceProperties map[string]interface{} `json:"OldResourceProperties,omitempty"`
}
