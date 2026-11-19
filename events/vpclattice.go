package events

// VPCLatticeRequestV1 contains a V1 request from AWS VPC Lattice.
type VPCLatticeRequestV1 struct {
	RawPath               string            `json:"raw_path"`
	Method                string            `json:"method"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"query_string_parameters"`
	Body                  string            `json:"body"`
	IsBase64Encoded       bool              `json:"is_base64_encoded"`
}

// VPCLatticeRequestV2 contains a V2 request from AWS VPC Lattice.
type VPCLatticeRequestV2 struct {
	Version               string                   `json:"version"`
	Path                  string                   `json:"path"`
	Method                string                   `json:"method"`
	Headers               map[string][]string      `json:"headers"`
	QueryStringParameters map[string][]string      `json:"queryStringParameters,omitempty"`
	Body                  string                   `json:"body"`
	RequestContext        VPCLatticeRequestContext `json:"requestContext"`
	IsBase64Encoded       bool                     `json:"isBase64Encoded,omitempty"`
}

// VPCLatticeRequestContext contains metadata about the incoming request.
type VPCLatticeRequestContext struct {
	ServiceNetworkARN string                     `json:"serviceNetworkArn"`
	ServiceARN        string                     `json:"serviceArn"`
	TargetGroupARN    string                     `json:"targetGroupArn"`
	Identity          *VPCLatticeRequestIdentity `json:"identity,omitempty"`
	Region            string                     `json:"region"`
	TimeEpoch         string                     `json:"timeEpoch"`
}

// VPCLatticeRequestIdentity contains information about the caller.
type VPCLatticeRequestIdentity struct {
	SourceVPCARN   string `json:"sourceVpcArn,omitempty"`
	Type           string `json:"type,omitempty"`
	Principal      string `json:"principal,omitempty"`
	PrincipalOrgID string `json:"principalOrgID,omitempty"`
	SessionName    string `json:"sessionName,omitempty"`
	X509IssuerOU   string `json:"x509IssuerOu,omitempty"`
	X509SanDNS     string `json:"x509SanDns,omitempty"`
	X509SanNameCN  string `json:"x509SanNameCn,omitempty"`
	X509SanURI     string `json:"x509SanUri,omitempty"`
	X509SubjectCN  string `json:"x509SubjectCn,omitempty"`
}

// VPCLatticeResponse contains the response to be returned to VPC Lattice.
type VPCLatticeResponse struct {
	IsBase64Encoded   bool              `json:"isBase64Encoded"`
	StatusCode        int               `json:"statusCode"`
	StatusDescription string            `json:"statusDescription,omitempty"`
	Headers           map[string]string `json:"headers"`
	Body              string            `json:"body,omitempty"`
}
