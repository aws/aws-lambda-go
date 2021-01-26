package events

// IoTCustomAuthorizerRequest contains data coming in to a custom IoT device gateway authorizer function.
type IoTCustomAuthorizerRequest struct {
	HttpContext        *IoTHttpContext `json:"httpContext,omitempty"`
	MqttContext        *IoTMqttContext `json:"mqttContext,omitempty"`
	TlsContext         *IotTlsContext  `json:"tlsContext,omitempty"`
	AuthorizationToken string          `json:"token"`
	TokenSignature     string          `json:"tokenSignature"`
}

type IoTHttpContext struct {
	Headers     map[string]string `json:"headers,omitempty"`
	QueryString string            `json:"queryString"`
}

type IoTMqttContext struct {
	ClientId string `json:"clientId"`
	Password []byte `json:"password"`
	Username string `json:"username"`
}

type IotTlsContext struct {
	ServerName string `json:"serverName"`
}

// IoTCustomAuthorizerResponse represents the expected format of an IoT device gateway authorization response.
type IoTCustomAuthorizerResponse struct {
	IsAuthenticated          bool                   `json:"isAuthenticated"`
	PrincipalID              string                 `json:"principalId"`
	DisconnectAfterInSeconds int32                  `json:"disconnectAfterInSeconds"`
	RefreshAfterInSeconds    int32                  `json:"refreshAfterInSeconds"`
	PolicyDocuments          []string               `json:"policyDocuments"`
}

// IoTCustomAuthorizerPolicy represents an IAM policy. PolicyDocuments is an array of IoTCustomAuthorizerPolicy JSON strings
type IoTCustomAuthorizerPolicy struct {
	Version   string
	Statement []IAMPolicyStatement
}
