package events

// IAMPolicyDocument represents an IAM policy document.
type IAMPolicyDocument struct {
	Version   string
	Statement []IAMPolicyStatement
}

// IAMPolicyStatement represents one statement from IAM policy with action, effect and resource.
type IAMPolicyStatement struct {
	Action   []string
	Effect   string
	Resource []string
	// Condition is an optional IAM condition block. Its value type is
	// interface{} because a condition value may be either a single string or an
	// array of strings, and omitempty keeps statements without conditions
	// serializing unchanged.
	Condition map[string]map[string]interface{} `json:"Condition,omitempty"`
}
