package events

// IoTCustomAuthorizerRequest contains data coming in to a custom IoT device gateway authorizer function.
type IoTCustomAuthorizerRequest struct {
	AuthorizationToken string `json:"token"`
}

// IoTCustomAuthorizerResponse represents the expected format of an IoT device gateway authorization response.
type IoTCustomAuthorizerResponse struct {
	IsAuthenticated          bool                        `json:"isAuthenticated"`
	PrincipalID              string                      `json:"principalId"`
	DisconnectAfterInSeconds int32                       `json:"disconnectAfterInSeconds"`
	RefreshAfterInSeconds    int32                       `json:"refreshAfterInSeconds"`
	PolicyDocuments          []IoTCustomAuthorizerPolicy `json:"policyDocuments"`
	Context                  map[string]interface{}      `json:"context,omitempty"`
}

// IoTCustomAuthorizerPolicy represents an IAM policy
type IoTCustomAuthorizerPolicy struct {
	Version   string
	Statement []IAMPolicyStatement
}

