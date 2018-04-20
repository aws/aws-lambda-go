// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

// CognitoEventPreSignup is sent by AWS Cognito when a user attempts to register
// (sign up), allowing a Lambda to perform custom validation to accept or deny the registration request
type CognitoEventPreSignup struct {
	CognitoEventHeader
	Request  CognitoEventPreSignupRequest  `json:"request"`
	Response CognitoEventPreSignupResponse `json:"response"`
}

// CognitoEventPostConfirmation is sent by AWS Cognito after a user is confirmed,
// allowing the Lambda to send custom messages or add custom logic.
type CognitoEventPostConfirmation struct {
	CognitoEventHeader
	Request  CognitoEventPostConfirmationRequest  `json:"request"`
	Response CognitoEventPostConfirmationResponse `json:"response"`
}

// CognitoEventCallerContext contains information about the caller
type CognitoEventCallerContext struct {
	AWSSDKVersion string `json:"awsSdkVersion"`
	ClientID      string `json:"clientId"`
}

// CognitoEventHeader contains common data from events sent by AWS Cognito
type CognitoEventHeader struct {
	Version       string                    `json:"version"`
	TriggerSource string                    `json:"triggerSource"`
	Region        string                    `json:"region"`
	UserPoolID    string                    `json:"userPoolId"`
	CallerContext CognitoEventCallerContext `json:"callerContext"`
	UserName      string                    `json:"userName"`
}

// CognitoEventPreSignupRequest contains the request portion of a PreSignup event
type CognitoEventPreSignupRequest struct {
	UserAttributes map[string]string `json:"userAttributes"`
	ValidationData map[string]string `json:"validationData"`
}

// CognitoEventPreSignupResponse contains the response portion of a PreSignup event
type CognitoEventPreSignupResponse struct {
	AutoConfirmUser bool `json:"autoConfirmUser"`
	AutoVerifyEmail bool `json:"autoVerifyEmail"`
	AutoVerifyPhone bool `json:"autoVerifyPhone"`
}

// CognitoEventPostConfirmationRequest contains the request portion of a PostConfirmation event
type CognitoEventPostConfirmationRequest struct {
	UserAttributes map[string]string `json:"userAttributes"`
}

// CognitoEventPostConfirmationResponse contains the response portion of a PostConfirmation event
type CognitoEventPostConfirmationResponse struct {
}
