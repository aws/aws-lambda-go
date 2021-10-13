package events

// IoTCustomAuthorizerRequest contains data coming in to a custom IoT device gateway authorizer function.
type IoTCustomAuthorizerRequest struct {
	Token              string                          `json:"token"`
	SignatureVerified  bool                            `json:"signatureVerified"` //whether the device gateway has validated the signature
	Protocols          []string                        `json:"protocols"`         //can include "tls", "mqtt", or "http"
	ProtocolData       IoTCustomAuthorizerProtocolData `json:"protocolData"`
	ConnectionMetadata IoTCustomAuthorizerMetadata     `json:"connectionMetadata"`
}

type IoTCustomAuthorizerProtocolData struct {
	HTTP *IoTHTTPContext `json:"http,omitempty"`
	MQTT *IoTMQTTContext `json:"mqtt,omitempty"`
	TLS  *IoTTLSContext  `json:"tls,omitempty"`
}

type IoTHTTPContext struct {
	Headers     map[string]string `json:"headers,omitempty"`
	QueryString string            `json:"queryString"`
}

type IoTMQTTContext struct {
	ClientID string `json:"clientId"`
	Password string `json:"password"` //base64-encoded string
	Username string `json:"username"`
}

type IoTTLSContext struct {
	ServerName string `json:"serverName"`
}

type IoTCustomAuthorizerMetadata struct {
	ID string `json:"id"` //UUID. The connection ID
}

// IoTCustomAuthorizerResponse represents the expected format of an IoT device gateway authorization response.
type IoTCustomAuthorizerResponse struct {
	IsAuthenticated          bool                               `json:"isAuthenticated"`
	PrincipalID              string                             `json:"principalId"`
	DisconnectAfterInSeconds int32                              `json:"disconnectAfterInSeconds"`
	RefreshAfterInSeconds    int32                              `json:"refreshAfterInSeconds"`
	PolicyDocuments          []APIGatewayCustomAuthorizerPolicy `json:"policyDocuments"`
}
