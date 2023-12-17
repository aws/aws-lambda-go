package events

import (
	"encoding/json"
	"time"
)

// EventBridgeEvent is the outer structure of an event sent via EventBridge serverless service.
// For examples of events that come via EventBridge, see https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-service-event.html
type EventBridgeEvent struct {
	Version    string          `json:"version"`
	ID         string          `json:"id"`
	DetailType string          `json:"detail-type"`
	Source     string          `json:"source"`
	AccountID  string          `json:"account"`
	Time       time.Time       `json:"time"`
	Region     string          `json:"region"`
	Resources  []string        `json:"resources"`
	Detail     json.RawMessage `json:"detail"`
}
