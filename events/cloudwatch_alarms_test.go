package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloudWatchAlarmEventMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
	}{
		{"metric alarm", "./testdata/cloudwatch-alarm-metric-event.json"},
		{"composite alarm", "./testdata/cloudwatch-alarm-composite-event.json"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inputJSON := test.ReadJSONFromFile(t, tc.filePath)

			var inputEvent CloudWatchAlarmEvent
			require.NoError(t, json.Unmarshal(inputJSON, &inputEvent))

			outputJSON, err := json.Marshal(inputEvent)
			require.NoError(t, err)

			assert.JSONEq(t, string(inputJSON), string(outputJSON))
		})
	}
}

func TestCloudWatchAlarmEventFields(t *testing.T) {
	var metricEvent CloudWatchAlarmEvent
	require.NoError(t, json.Unmarshal(test.ReadJSONFromFile(t, "./testdata/cloudwatch-alarm-metric-event.json"), &metricEvent))

	assert.Equal(t, "ALARM", metricEvent.AlarmData.State.Value)
	assert.Equal(t, "INSUFFICIENT_DATA", metricEvent.AlarmData.PreviousState.Value)
	require.Len(t, metricEvent.AlarmData.Configuration.Metrics, 1)
	metric := metricEvent.AlarmData.Configuration.Metrics[0]
	require.NotNil(t, metric.MetricStat)
	assert.Equal(t, "i-12345678", metric.MetricStat.Metric.Dimensions["InstanceId"])
	assert.True(t, metric.ReturnData)
	assert.Equal(t, []CloudWatchAlarmTag{{Key: "Environment", Value: "prod"}}, metricEvent.AlarmData.AlarmTags)

	var compositeEvent CloudWatchAlarmEvent
	require.NoError(t, json.Unmarshal(test.ReadJSONFromFile(t, "./testdata/cloudwatch-alarm-composite-event.json"), &compositeEvent))

	assert.Equal(t, "ALARM(CompositeDemo.FirstChild) OR ALARM(CompositeDemo.SecondChild)", compositeEvent.AlarmData.Configuration.AlarmRule)
	assert.Equal(t, int64(120), compositeEvent.AlarmData.Configuration.ActionsSuppressorWaitPeriod)
	assert.Equal(t, "WaitPeriod", compositeEvent.AlarmData.PreviousState.ActionsSuppressedBy)
	assert.Empty(t, compositeEvent.AlarmData.Configuration.Metrics)
}

func TestCloudWatchAlarmEventMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, &CloudWatchAlarmEvent{})
}
