package events

// LambdaTargetGroupRequest contains data originating from the ALB Lambda target group integration
type LambdaTargetGroupRequest struct {
	HTTPMethod                      string                          `json:"httpMethod"`
	Path                            string                          `json:"path"`
	QueryStringParameters           map[string]string               `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string             `json:"multiValueQueryStringParameters"`
	Headers                         map[string]string               `json:"headers"`
	MultiValueHeaders               map[string][]string             `json:"multiValueHeaders"`
	RequestContext                  LambdaTargetGroupRequestContext `json:"requestContext"`
	IsBase64Encoded                 bool                            `json:"isBase64Encoded"`
	Body                            string                          `json:"body"`
}

// LambdaTargetGroupRequestContext contains the information to identify the load balancer invoking the lambda
type LambdaTargetGroupRequestContext struct {
	Elb ELBContext `json:"elb"`
}

// ELBContext contains the information to identify the ARN invoking the lambda
type ELBContext struct {
	TargetGroupArn string `json:"targetGroupArn"`
}

// LambdaTargetGroupResponse configures the response to be returned by the ALB Lambda target group for the request
type LambdaTargetGroupResponse struct {
	StatusCode        int                 `json:"statusCode"`
	StatusDescription string              `json:"statusDescription"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded"`
}
