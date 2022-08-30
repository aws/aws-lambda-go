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

//AppSyncDirectLambdaResolverRequest represents a direct lambda resolver request event
type AppSyncDirectLambdaResolverRequest struct {
	Arguments map[string]interface{} `json:"arguments"`
	Identity  map[string]interface{} `json:"identity"`
	Info      struct {
		FieldName           string   `json:"fieldName"`
		ParentTypeName      string   `json:"parentTypeName"`
		SelectionSetGraphQL string   `json:"selectionSetGraphQL"`
		SelectionSetList    []string `json:"selectionSetList"`
	}
	Prev    string `json:"prev"`
	Request struct {
		Headers struct {
			Accept                    string `json:"accept"`
			Authorization             string `json:"authorization"`
			AcceptEncoding            string `json:"accept-encoding"`
			AcceptLanguage            string `json:"accept-language"`
			CloudfrontForwardedProto  string `json:"cloudfront-forwarded-proto"`
			CloudfrontIsDesktopViewer string `json:"cloudfront-is-desktop-viewer"`
			CloudfrontIsMobileViewer  string `json:"cloudfront-is-mobile-viewer"`
			CloudfrontIsSmarttvViewer string `json:"cloudfront-is-smarttv-viewer"`
			CloudfrontViewerCountry   string `json:"cloudfront-viewer-country"`
			CloudfrontIsTabletViewer  string `json:"cloudfront-is-tablet-viewer"`
			ContentLength             string `json:"content-length"`
			ContentType               string `json:"content-type"`
			Host                      string `json:"host"`
			Hrigin                    string `json:"origin"`
			Referer                   string `json:"Referer"`
			SecFetchDest              string `json:"sec-fetch-dest"`
			SecFetchMode              string `json:"sec-fetch-mode"`
			SecFetchSite              string `json:"sec-fetch-site"`
			UserAgent                 string `json:"user-agent"`
			Via                       string `json:"via"`
			XAmzCfID                  string `json:"x-amz-cf-id"`
			XAmzUserAgent             string `json:"x-amz-user-agent"`
			XAmznTraceID              string `json:"x-amzn-trace-id"`
			XApiKey                   string `json:"x-api-key"`
			XForwardedFor             string `json:"x-forwarded-for"`
			XForwardedPort            string `json:"x-forwarded-port"`
			XForwardedProto           string `json:"x-forwarded-proto"`
		}
	}
	Source map[string]interface{} `json:"source"`
	Stash  map[string]string      `json:"stash"`
}
