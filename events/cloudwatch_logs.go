package events

// CloudwatchLogsEvent represents raw data from a cloudwatch logs event
type CloudwatchLogsEvent struct {
	AWSLogs CloudwatchLogsRawData `json:"awslogs"`
}

// CloudwatchLogsRawData contains gzipped base64 json representing the bulk
// of a cloudwatch logs event
type CloudwatchLogsRawData struct {
	Data string `json:"data"`
}
