// Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

// LambdaHTTPRequest contains data coming from the new HTTP API Gateway
type LambdaHTTPRequest struct {
	Version               string                   `json:"version"`  // Version is expected to always be `"2.0"`
	RouteKey              string                   `json:"routeKey"` // RouteKey is expected to always be `"$default"`
	RawPath               string                   `json:"rawPath"`
	RawQueryString        string                   `json:"rawQueryString"`
	Cookies               []string                 `json:"cookies,omitempty"`
	Headers               map[string]string        `json:"headers"`
	QueryStringParameters map[string]string        `json:"queryStringParameters,omitempty"`
	RequestContext        LambdaHTTPRequestContext `json:"requestContext"`
	Body                  string                   `json:"body,omitempty"`
	IsBase64Encoded       bool                     `json:"isBase64Encoded"`
}

// LambdaHTTPRequestContext contains the information to identify the AWS account and resources invoking the Lambda function.
type LambdaHTTPRequestContext struct {
	RouteKey     string                                         `json:"routeKey"` // RouteKey is expected to always be `"$default"`
	AccountID    string                                         `json:"accountId"`
	Stage        string                                         `json:"stage"` // Stage is expected to always be `"$default"`
	RequestID    string                                         `json:"requestId"`
	Authorizer   *LambdaHTTPRequestContextAuthorizerDescription `json:"authorizer,omitempty"`
	APIID        string                                         `json:"apiId"`        // APIID is the Lambda URL ID
	DomainName   string                                         `json:"domainName"`   // DomainName is of the format `"<url-id>.lambda-url.<region>.on.aws"`
	DomainPrefix string                                         `json:"domainPrefix"` // DomainPrefix is the Lambda URL ID
	Time         string                                         `json:"time"`
	TimeEpoch    int64                                          `json:"timeEpoch"`
	HTTP         LambdaHTTPRequestContextHTTPDescription        `json:"http"`
}

// LambdaHTTPRequestContextAuthorizerDescription contains authorizer information for the request context.
type LambdaHTTPRequestContextAuthorizerDescription struct {
	IAM *LambdaHTTPRequestContextAuthorizerIAMDescription `json:"iam,omitempty"`
}

// LambdaHTTPRequestContextAuthorizerIAMDescription contains IAM information for the request context.
type LambdaHTTPRequestContextAuthorizerIAMDescription struct {
	AccessKey string `json:"accessKey"`
	AccountID string `json:"accountId"`
	CallerID  string `json:"callerId"`
	UserARN   string `json:"userArn"`
	UserID    string `json:"userId"`
}

// LambdaHTTPRequestContextHTTPDescription contains HTTP information for the request context.
type LambdaHTTPRequestContextHTTPDescription struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	SourceIP  string `json:"sourceIp"`
	UserAgent string `json:"userAgent"`
}

// LambdaHTTPResponse configures the response to be returned by API Gateway V2 for the request
type LambdaHTTPResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	Cookies         []string          `json:"cookies"`
}
