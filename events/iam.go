package events

type IAMPolicyStatement struct {
	Action   []string
	Effect   string
	Resource []string
}
