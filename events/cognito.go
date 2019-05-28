// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

// CognitoEvent contains data from an event sent from AWS Cognito Sync
type CognitoEvent struct {
	DatasetName    string                          `json:"datasetName"`
	DatasetRecords map[string]CognitoDatasetRecord `json:"datasetRecords"`
	EventType      string                          `json:"eventType"`
	IdentityID     string                          `json:"identityId"`
	IdentityPoolID string                          `json:"identityPoolId"`
	Region         string                          `json:"region"`
	Version        int                             `json:"version"`
}

// CognitoDatasetRecord represents a record from an AWS Cognito Sync event
type CognitoDatasetRecord struct {
	NewValue string `json:"newValue"`
	OldValue string `json:"oldValue"`
	Op       string `json:"op"`
}

// CognitoEventUserPoolsPreSignup is sent by AWS Cognito User Pools when a user attempts to register
// (sign up), allowing a Lambda to perform custom validation to accept or deny the registration request
type CognitoEventUserPoolsPreSignup struct {
	CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsPreSignupRequest  `json:"request"`
	Response CognitoEventUserPoolsPreSignupResponse `json:"response"`
}

// CognitoEventUserPoolsPostConfirmation is sent by AWS Cognito User Pools after a user is confirmed,
// allowing the Lambda to send custom messages or add custom logic.
type CognitoEventUserPoolsPostConfirmation struct {
	CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsPostConfirmationRequest  `json:"request"`
	Response CognitoEventUserPoolsPostConfirmationResponse `json:"response"`
}

// CognitoEventUserPoolsPreTokenGen is sent by AWS Cognito User Pools when a user attempts to retrieve
// credentials, allowing a Lambda to perform insert, supress or override claims
type CognitoEventUserPoolsPreTokenGen struct {
	CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsPreTokenGenRequest  `json:"request"`
	Response CognitoEventUserPoolsPreTokenGenResponse `json:"response"`
}

// CognitoEventUserPoolsPostAuthentication is sent by AWS Cognito User Pools after a user is authenticated,
// allowing the Lambda to add custom logic.
type CognitoEventUserPoolsPostAuthentication struct {
	CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsPostAuthenticationRequest  `json:"request"`
	Response CognitoEventUserPoolsPostAuthenticationResponse `json:"response"`
}

// CognitoEventUserPoolsMigrateUser is sent by AWS Cognito User Pools when a user does not exist in the
// user pool at the time of sign-in with a password, or in the forgot-password flow.
type CognitoEventUserPoolsMigrateUser struct {
	CognitoEventUserPoolsHeader
	CognitoEventUserPoolsMigrateUserRequest  `json:"request"`
	CognitoEventUserPoolsMigrateUserResponse `json:"response"`
}

// CognitoEventUserPoolsCallerContext contains information about the caller
type CognitoEventUserPoolsCallerContext struct {
	AWSSDKVersion string `json:"awsSdkVersion"`
	ClientID      string `json:"clientId"`
}

// CognitoEventUserPoolsHeader contains common data from events sent by AWS Cognito User Pools
type CognitoEventUserPoolsHeader struct {
	Version       string                             `json:"version"`
	TriggerSource string                             `json:"triggerSource"`
	Region        string                             `json:"region"`
	UserPoolID    string                             `json:"userPoolId"`
	CallerContext CognitoEventUserPoolsCallerContext `json:"callerContext"`
	UserName      string                             `json:"userName"`
}

// CognitoEventUserPoolsPreSignupRequest contains the request portion of a PreSignup event
type CognitoEventUserPoolsPreSignupRequest struct {
	UserAttributes map[string]string `json:"userAttributes"`
	ValidationData map[string]string `json:"validationData"`
}

// CognitoEventUserPoolsPreSignupResponse contains the response portion of a PreSignup event
type CognitoEventUserPoolsPreSignupResponse struct {
	AutoConfirmUser bool `json:"autoConfirmUser"`
	AutoVerifyEmail bool `json:"autoVerifyEmail"`
	AutoVerifyPhone bool `json:"autoVerifyPhone"`
}

// CognitoEventUserPoolsPostConfirmationRequest contains the request portion of a PostConfirmation event
type CognitoEventUserPoolsPostConfirmationRequest struct {
	UserAttributes map[string]string `json:"userAttributes"`
}

// CognitoEventUserPoolsPostConfirmationResponse contains the response portion of a PostConfirmation event
type CognitoEventUserPoolsPostConfirmationResponse struct {
}

// CognitoEventUserPoolsPreTokenGenRequest contains request portion of PreTokenGen event
type CognitoEventUserPoolsPreTokenGenRequest struct {
	UserAttributes     map[string]string  `json:"userAttributes"`
	GroupConfiguration GroupConfiguration `json:"groupConfiguration"`
}

// CognitoEventUserPoolsPreTokenGenResponse containst the response portion of  a PreTokenGen event
type CognitoEventUserPoolsPreTokenGenResponse struct {
	ClaimsOverrideDetails ClaimsOverrideDetails `json:"claimsOverrideDetails"`
}

// CognitoEventUserPoolsPostAuthenticationRequest contains the request portion of a PostAuthentication event
type CognitoEventUserPoolsPostAuthenticationRequest struct {
	NewDeviceUsed  bool              `json:"newDeviceUsed"`
	UserAttributes map[string]string `json:"userAttributes"`
}

// CognitoEventUserPoolsPostAuthenticationResponse contains the response portion of a PostAuthentication event
type CognitoEventUserPoolsPostAuthenticationResponse struct {
}

// CognitoEventUserPoolsMigrateUserRequest contains the request portion of a MigrateUser event
type CognitoEventUserPoolsMigrateUserRequest struct {
	Password string `json:"password"`
}

// CognitoEventUserPoolsMigrateUserResponse contains the response portion of a MigrateUser event
type CognitoEventUserPoolsMigrateUserResponse struct {
	UserAttributes         map[string]string `json:"userAttributes"`
	FinalUserStatus        string            `json:"finalUserStatus"`
	MessageAction          string            `json:"messageAction"`
	DesiredDeliveryMediums []string          `json:"desiredDeliveryMediums"`
	ForceAliasCreation     bool              `json:"forceAliasCreation"`
}

// ClaimsOverrideDetails allows lambda to add, supress or override claims in the token
type ClaimsOverrideDetails struct {
	GroupOverrideDetails  GroupConfiguration `json:"groupOverrideDetails"`
	ClaimsToAddOrOverride map[string]string  `json:"claimsToAddOrOverride"`
	ClaimsToSuppress      []string           `json:"claimsToSuppress"`
}

// GroupConfiguration allows lambda to override groups, roles and set a perferred role
type GroupConfiguration struct {
	GroupsToOverride   []string `json:"groupsToOverride"`
	IAMRolesToOverride []string `json:"iamRolesToOverride"`
	PreferredRole      *string  `json:"preferredRole"`
}

// CognitoEventUserPoolsCustomMessage is sent by AWS Cognito User Pools before a verification or MFA message is sent,
// allowing a user to customize the message dynamically.
type CognitoEventUserPoolsCustomMessage struct {
	CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsCustomMessageRequest  `json:"request"`
	Response CognitoEventUserPoolsCustomMessageResponse `json:"response"`
}

// CognitoEventUserPoolsCustomMessageRequest contains the request portion of a CustomMessage event
type CognitoEventUserPoolsCustomMessageRequest struct {
	UserAttributes map[string]interface{} `json:"userAttributes"`
	CodeParameter  string                 `json:"codeParameter"`
}

// CognitoEventUserPoolsCustomMessageResponse contains the response portion of a CustomMessage event
type CognitoEventUserPoolsCustomMessageResponse struct {
	SMSMessage   string `json:"smsMessage"`
	EmailMessage string `json:"emailMessage"`
	EmailSubject string `json:"emailSubject"`
}
