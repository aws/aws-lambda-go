package events

import "encoding/json"

// Deprecated: AppSyncResolverTemplate does not represent resolver events sent by AppSync. Instead directly model your input schema, or use map[string]string, json.RawMessage, interface{}, etc..
type AppSyncResolverTemplate struct {
	Version   string           `json:"version"`
	Operation AppSyncOperation `json:"operation"`
	Payload   json.RawMessage  `json:"payload"`
}

// AppSyncIAMIdentity contains information about the caller authed via IAM.
type AppSyncIAMIdentity struct {
	AccountID                   string   `json:"accountId"`
	CognitoIdentityAuthProvider string   `json:"cognitoIdentityAuthProvider"`
	CognitoIdentityAuthType     string   `json:"cognitoIdentityAuthType"`
	CognitoIdentityPoolID       string   `json:"cognitoIdentityPoolId"`
	CognitoIdentityID           string   `json:"cognitoIdentityId"`
	SourceIP                    []string `json:"sourceIp"`
	Username                    string   `json:"username"`
	UserARN                     string   `json:"userArn"`
}

// AppSyncCognitoIdentity contains information about the caller authed via Cognito.
type AppSyncCognitoIdentity struct {
	Sub                 string                 `json:"sub"`
	Issuer              string                 `json:"issuer"`
	Username            string                 `json:"username"`
	Claims              map[string]interface{} `json:"claims"`
	SourceIP            []string               `json:"sourceIp"`
	DefaultAuthStrategy string                 `json:"defaultAuthStrategy"`
}

// Deprecated: not used by any event schema
type AppSyncOperation string

const (
	// Deprecated: not used by any event schema
	OperationInvoke AppSyncOperation = "Invoke"
	// Deprecated: not used by any event schema
	OperationBatchInvoke AppSyncOperation = "BatchInvoke"
)

// AppSyncLambdaAuthorizerRequest contains an authorization request from AppSync.
type AppSyncLambdaAuthorizerRequest struct {
	AuthorizationToken string                                `json:"authorizationToken"`
	RequestContext     AppSyncLambdaAuthorizerRequestContext `json:"requestContext"`
}

// AppSyncLambdaAuthorizerRequestContext contains the parameters of the AppSync invocation which triggered
// this authorization request.
type AppSyncLambdaAuthorizerRequestContext struct {
	APIID         string                 `json:"apiId"`
	AccountID     string                 `json:"accountId"`
	RequestID     string                 `json:"requestId"`
	QueryString   string                 `json:"queryString"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// AppSyncLambdaAuthorizerResponse represents the expected format of an authorization response to AppSync.
type AppSyncLambdaAuthorizerResponse struct {
	IsAuthorized    bool                   `json:"isAuthorized"`
	ResolverContext map[string]interface{} `json:"resolverContext,omitempty"`
	DeniedFields    []string               `json:"deniedFields,omitempty"`
	TTLOverride     *int                   `json:"ttlOverride,omitempty"`
}

const (
	// AllowAuth is allowing by default.
	AuthStrategy_Allow AuthStrategy = "ALLOW"
	// DenyAuth is denying by default.
	AuthStrategy_Deny AuthStrategy = "DENY"
)

type (
	AuthStrategy string

	AppSyncDirectLambdaResolverRequest struct {
		Arguments map[string]interface{} `json:"arguments"`
		Source    map[string]interface{} `json:"source"`
		Result    interface{}            `json:"result"`
		Identity  Identity               `json:"identity"`
		Request   Request                `json:"request"`
		Info      Info                   `json:"info"`
	}

	Info struct {
		FieldName           string                 `json:"fieldName"`
		ParentTypeName      string                 `json:"parentTypeName"`
		SelectionSetGraphQL string                 `json:"selectionSetGraphQL"`
		Variables           map[string]interface{} `json:"variables"`
		SelectionSetList    []string               `json:"selectionSetList"`
	}

	Request struct {
		DomainName *string           `json:"domainName"`
		Headers    map[string]string `json:"headers"`
	}

	Identity struct {
		// User Pool
		Sub                 string                 `json:"sub"`
		Issuer              string                 `json:"issuer"`
		Username            string                 `json:"username"`
		DefaultAuthStrategy AuthStrategy           `json:"defaultAuthStrategy"`
		Claims              map[string]interface{} `json:"claims"`
		// AWS IAM
		CognitoIdentityID           string `json:"cognitoIdentityId"`
		CognitoIdentityPoolID       string `json:"cognitoIdentityPoolId"`
		CognitoIdentityAuthType     string `json:"cognitoIdentityAuthType"`
		AccountID                   string `json:"accountId"`
		CognitoIdentityAuthProvider string `json:"cognitoIdentityAuthProvider"`
		// common User Pool and AWS_IAM
		SourceIP []string `json:"sourceIp"`
	}
)
