package events

type CustomResourceRequestType = string

// Enum values for CustomResourceRequestType
const (
	CustomResourceRequestCreate CustomResourceRequestType = "Create"
	CustomResourceRequestUpdate CustomResourceRequestType = "Update"
	CustomResourceRequestDelete CustomResourceRequestType = "Delete"
)

// CloudformationCustomResourceRequest represents a Custom Resource Request
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/crpg-ref-requests.html
type CloudformationCustomResourceRequest struct {
	RequestType           CustomResourceRequestType
	ResponseURL           string
	StackId               string
	RequestId             string
	ResourceType          string
	LogicalResourceId     string
	PhysicalResourceId    string `json:",omitempty"`
	ResourceProperties    map[string]interface{}
	OldResourceProperties map[string]interface{} `json:",omitempty"`
}

type CustomResourceStatus = string

// Enum values for CustomResourceStatus
const (
	CustomResourceSuccess CustomResourceStatus = "SUCCESS"
	CustomResourceFailed  CustomResourceStatus = "FAILED"
)

// CloudformationCustomResourceResponse represents a Custom Resource Response
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/crpg-ref-responses.html
type CloudformationCustomResourceResponse struct {
	Status             CustomResourceStatus
	Reason             string `json:",omitempty"`
	PhysicalResourceId string
	StackId            string
	RequestId          string
	LogicalResourceId  string
	NoEcho             bool `json:",omitempty"`
	Data               map[string]interface{}
}
