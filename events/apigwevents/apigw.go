// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package apigwevents

// ApiGatewayProxyRequest contains data coming from the API Gateway proxy
type ApiGatewayProxyRequest struct {
	Resource              string              `json:"resource"` // The resource path defined in API Gateway
	Path                  string              `json:"path"`     // The url path for the caller
	HttpMethod            string              `json:"httpMethod"`
	Headers               map[string]string   `json:"headers"`
	QueryStringParameters map[string]string   `json:"queryStringParameters"`
	PathParameters        map[string]string   `json:"pathParameters"`
	StageVariables        map[string]string   `json:"stageVariables"`
	RequestContext        ProxyRequestContext `json:"requestContext"`
	Body                  string              `json:"body"`
	IsBase64Encoded       bool                `json:"isBase64Encoded,omitempty"`
}

type ApiGatewayProxyResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded,omitempty"`
}

// ProxyRequestContext contains the information to identify the AWS account and resources invoking the
// Lambda function. It also includes Cognito identity information for the caller.
type ProxyRequestContext struct {
	AccountId    string          `json:"accountId"`
	ResourceId   string          `json:"resourceId"`
	Stage        string          `json:"stage"`
	RequestId    string          `json:"requestId"`
	Identity     RequestIdentity `json:"identity"`
	ResourcePath string          `json:"resourcePath"`
	HttpMethod   string          `json:"httpMethod"`
	ApiId        string          `json:"apiId"` // The API Gateway rest API Id
}

// RequestIdentity contains identity information for the request caller.
type RequestIdentity struct {
	CognitoIdentityPoolId         string `json:"cognitoIdentityPoolId"`
	AccountId                     string `json:"accountId"`
	CognitoIdentityId             string `json:"cognitoIdentityId"`
	Caller                        string `json:"caller"`
	ApiKey                        string `json:"apiKey"`
	SourceIp                      string `json:"sourceIp"`
	CognitoAuthenticationType     string `json:"cognitoAuthenticationType"`
	CognitoAuthenticationProvider string `json:"cognitoAuthenticationProvider"`
	UserArn                       string `json:"userArn"`
	UserAgent                     string `json:"userAgent"`
	User                          string `json:"user"`
}

// ApiGatewayCustomAuthorizerContext represents the expected format of an API Gateway custom authorizer response.
type ApiGatewayCustomAuthorizerContext struct {
	PrincipalId *string `json:"principalId"`
	StringKey   *string `json:"stringKey,omitempty"`
	NumKey      *int    `json:"numKey,omitempty"`
	BoolKey     *bool   `json:"boolKey,omitempty"`
}
