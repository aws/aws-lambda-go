package events

import "encoding/json"

// AppSyncResolverTemplate represents the requests from AppSync to Lambda
type AppSyncResolverTemplate struct {
	Version   string           `json:"version"`
	Operation AppSyncOperation `json:"operation"`
	Payload   json.RawMessage  `json:"payload"`
}

// AppSyncIdentity contains information about the caller. The shape of this section depends on the authorization type of your AWS AppSync API
type AppSyncIdentity struct {
	AccountId             string                 `json:"accountId"`
	Claims                map[string]interface{} `json:"claims"`
	CognitoIdentityPoolId string                 `json:"cognitoIdentityPoolId"`
	CognitoIdentityId     string                 `json:"cognitoIdentityId"`
	DefaultAuthStrategy   string                 `json:"defaultAuthStrategy"`
	Issuer                string                 `json:"issuer"`
	SourceIp              []string               `json:"sourceIp"`
	Sub                   string                 `json:"uuid"`
	Username              string                 `json:"username"`
	UserArn               string                 `json:"userArn"`
}

// AppSyncOperation specifies the operation type supported by Lambda operations
type AppSyncOperation string

const (
	// OperationInvoke lets AWS AppSync know to call your Lambda function for every GraphQL field resolver
	OperationInvoke AppSyncOperation = "Invoke"
	// OperationBatchInvoke instructs AWS AppSync to batch requests for the current GraphQL field
	OperationBatchInvoke AppSyncOperation = "BatchInvoke"
)
