package events

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
)

// EventBridgeEventLogs represents raw data from an eventbridge logs event
type EventBridgeEventLogs struct {
	AWSLogs EventBridgeLogsRawData `json:"awslogs"`
}

// EventBridgeLogsRawData contains gzipped base64 json representing the bulk
// of an eventbridge logs event
type EventBridgeLogsRawData struct {
	Data string `json:"data"`
}

// Parse returns a struct representing a usable EventBridgeLogs event
func (c EventBridgeLogsRawData) Parse() (d EventBridgeLogsData, err error) {
	data, err := base64.StdEncoding.DecodeString(c.Data)
	if err != nil {
		return
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	defer zr.Close()

	dec := json.NewDecoder(zr)
	err = dec.Decode(&d)

	return
}

// EventBridgeLogsData is an unmarshal'd, ungzip'd, eventbridge logs event
type EventBridgeLogsData struct {
	Owner               string                `json:"owner"`
	LogGroup            string                `json:"logGroup"`
	LogStream           string                `json:"logStream"`
	SubscriptionFilters []string              `json:"subscriptionFilters"`
	MessageType         string                `json:"messageType"`
	LogEvents           []EventBridgeLogEvent `json:"logEvents"`
}

// EventBridgeLogEvent represents a log entry from eventbridge logs
type EventBridgeLogEvent struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}
