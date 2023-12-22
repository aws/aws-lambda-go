package events

type S3ObjectLambdaEvent struct {
	XAmzRequestId        string                      `json:"xAmzRequestId"`
	GetObjectContext     *GetObjectContext           `json:"getObjectContext,omitempty"`
	ListObjectsContext   *ListObjectsContext         `json:"listObjectsContext,omitempty"`
	ListObjectsV2Context *ListObjectsV2Context       `json:"listObjectsV2Context,omitempty"`
	HeadObjectContext    *HeadObjectContext          `json:"headObjectContext,omitempty"`
	Configuration        S3ObjectLambdaConfiguration `json:"configuration"`
	UserRequest          UserRequest                 `json:"userRequest"`
	UserIdentity         UserIdentity                `json:"userIdentity"`
	ProtocolVersion      string                      `json:"protocolVersion"`
}

type GetObjectContext struct {
	InputS3Url  string `json:"inputS3Url"`
	OutputRoute string `json:"outputRoute"`
	OutputToken string `json:"outputToken"`
}

type ListObjectsContext struct {
	InputS3Url string `json:"inputS3Url"`
}

type ListObjectsV2Context struct {
	InputS3Url string `json:"inputS3Url"`
}

type HeadObjectContext struct {
	InputS3Url string `json:"inputS3Url"`
}

type S3ObjectLambdaConfiguration struct {
	AccessPointArn           string `json:"accessPointArn"`
	SupportingAccessPointArn string `json:"supportingAccessPointArn"`
	Payload                  string `json:"payload"`
}

type UserRequest struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type UserIdentity struct {
	Type           string          `json:"type"`
	PrincipalId    string          `json:"principalId"`
	Arn            string          `json:"arn"`
	AccountId      string          `json:"accountId"`
	AccessKeyId    string          `json:"accessKeyId"`
	SessionContext *SessionContext `json:"sessionContext,omitempty"`
}

type SessionContext struct {
	Attributes    map[string]string `json:"attributes"`
	SessionIssuer *SessionIssuer    `json:"sessionIssuer,omitempty"`
}

type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalId string `json:"principalId"`
	Arn         string `json:"arn"`
	AccountId   string `json:"accountId"`
	UserName    string `json:"userName"`
}
