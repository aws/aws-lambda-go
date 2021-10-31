package events

type CustomResourceEvent struct {
	RequestType           string                 `json:"RequestType"` // One of Create, Update, or Delete.
	ResponseURL           string                 `json:"ResponseURL"`
	StackID               string                 `json:"StackId"`
	RequestID             string                 `json:"RequestId"`
	ResourceType          string                 `json:"ResourceType"`
	LogicalResourceID     string                 `json:"LogicalResourceId"`
	PhysicalResourceID    string                 `json:"PhysicalResourceId,omitempty"` // Always sent with Update and Delete requests; never sent with Create.
	ResourceProperties    map[string]interface{} `json:"ResourceProperties"`
	OldResourceProperties map[string]interface{} `json:"OldResourceProperties,omitempty"` // Used only for Update requests.
}
