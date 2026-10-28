package events

// CloudWatchAlarmEvent is the event sent to a Lambda function when it is configured as a CloudWatch alarm action.
// See https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/alarms-and-actions-Lambda.html
type CloudWatchAlarmEvent struct {
	Source    string              `json:"source"`
	AlarmARN  string              `json:"alarmArn"`
	AccountID string              `json:"accountId"`
	Time      string              `json:"time"`
	Region    string              `json:"region"`
	AlarmData CloudWatchAlarmData `json:"alarmData"`
}

// CloudWatchAlarmData describes the alarm, its current and previous state, and its configuration.
type CloudWatchAlarmData struct {
	AlarmName       string                       `json:"alarmName"`
	State           CloudWatchAlarmState         `json:"state"`
	PreviousState   CloudWatchAlarmState         `json:"previousState"`
	Configuration   CloudWatchAlarmConfiguration `json:"configuration"`
	AlarmTags       []CloudWatchAlarmTag         `json:"alarmTags,omitempty"`
	AlarmTagsStatus string                       `json:"alarmTagsStatus,omitempty"`
}

// CloudWatchAlarmState is a state of the alarm. ReasonData is a JSON-encoded string.
type CloudWatchAlarmState struct {
	Value                   string `json:"value"`
	Reason                  string `json:"reason"`
	ReasonData              string `json:"reasonData,omitempty"`
	Timestamp               string `json:"timestamp"`
	ActionsSuppressedBy     string `json:"actionsSuppressedBy,omitempty"`
	ActionsSuppressedReason string `json:"actionsSuppressedReason,omitempty"`
}

// CloudWatchAlarmConfiguration holds the configuration of a metric alarm (Metrics)
// or a composite alarm (AlarmRule and the ActionsSuppressor fields).
type CloudWatchAlarmConfiguration struct {
	Description                      string                           `json:"description,omitempty"`
	Metrics                          []CloudWatchAlarmMetricDataQuery `json:"metrics,omitempty"`
	AlarmRule                        string                           `json:"alarmRule,omitempty"`
	ActionsSuppressor                string                           `json:"actionsSuppressor,omitempty"`
	ActionsSuppressorWaitPeriod      int64                            `json:"actionsSuppressorWaitPeriod,omitempty"`
	ActionsSuppressorExtensionPeriod int64                            `json:"actionsSuppressorExtensionPeriod,omitempty"`
}

type CloudWatchAlarmMetricDataQuery struct {
	ID         string                     `json:"id"`
	Expression string                     `json:"expression,omitempty"`
	Label      string                     `json:"label,omitempty"`
	MetricStat *CloudWatchAlarmMetricStat `json:"metricStat,omitempty"`
	ReturnData bool                       `json:"returnData"`
}

type CloudWatchAlarmMetricStat struct {
	Metric CloudWatchAlarmMetric `json:"metric"`
	Period int64                 `json:"period"`
	Stat   string                `json:"stat"`
	Unit   string                `json:"unit,omitempty"`
}

type CloudWatchAlarmMetric struct {
	Namespace  string            `json:"namespace"`
	Name       string            `json:"name"`
	Dimensions map[string]string `json:"dimensions,omitempty"`
}

type CloudWatchAlarmTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
