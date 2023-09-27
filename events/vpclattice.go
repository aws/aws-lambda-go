package events

// VPCLatticeRequestV1 contains data coming from AWS VPC Lattice
type VPCLatticeRequestV1 struct {
	RawPath               string            `json:"raw_path"`
	Method                string            `json:"method"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"query_string_parameters"`
	Body                  string            `json:"body"`
	IsBase64Encoded       bool              `json:"is_base64_encoded,omitempty"`
}

// VPCLatticeRequestV2 contains data coming from AWS VPC Lattice
type VPCLatticeRequestV2 struct {
	Version        string                   `json:"version"`
	Path           string                   `json:"path"`
	Method         string                   `json:"method"`
	Headers        map[string][]string      `json:"headers"`
	Body           string                   `json:"body"`
	RequestContext VpcLatticeRequestContext `json:"requestContext"`
}

// VpcLatticeRequestContext contains metadata about the incoming request
type VpcLatticeRequestContext struct {
	ServiceNetworkArn string                     `json:"serviceNetworkArn"`
	ServiceArn        string                     `json:"serviceArn"`
	TargetGroupArn    string                     `json:"targetGroupArn"`
	Identity          *VpcLatticeRequestIdentity `json:"identity,omitempty"`
	Region            string                     `json:"region"`
	TimeEpoch         string                     `json:"timeEpoch"`
}

// VpcLatticeRequestIdentity contains information about the caller
type VpcLatticeRequestIdentity struct {
	SourceVpcArn string `json:"sourceVpcArn"`
	Type         string `json:"type"`
	Principal    string `json:"principal"`
	SessionName  string `json:"sessionName"`
}

// VPCLatticeResponse contains the response to be returned to VPC Lattice
type VPCLatticeResponse struct {
	IsBase64Encoded   bool              `json:"isBase64Encoded,omitempty"`
	StatusCode        int               `json:"statusCode"`
	StatusDescription string            `json:"statusDescription,omitempty"`
	Headers           map[string]string `json:"headers"`
	Body              string            `json:"body"`
}