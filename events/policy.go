package events

// CustomAuthorizerPolicy represents an IAM policy
type CustomAuthorizerPolicy struct {
	Version   string
	Statement []IAMPolicyStatement
}
