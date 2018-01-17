package events

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCloudwatchLogsParse(t *testing.T) {
	for _, test := range []struct {
		name                     string
		eventJson                string
		expectError              bool
		expectCloudwatchLogsData CloudwatchLogsData
	}{
		{"Well formed cloudwatch event",
			"./testdata/cloudwatch-logs-event.json",
			false,
			CloudwatchLogsData{
				Owner:     "123456789123",
				LogGroup:  "testLogGroup",
				LogStream: "testLogStream",
				SubscriptionFilters: []string{
					"testFilter",
				},
				MessageType: "DATA_MESSAGE",
				LogEvents: []LogEvent{
					{ID: "eventId1", Timestamp: 1440442987000, Message: "[ERROR] First test message"},
					{ID: "eventId2", Timestamp: 1440442987001, Message: "[ERROR], Second test message"},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			inputJson := readJsonFromFile(t, test.eventJson)

			var inputEvent CloudwatchLogsEvent
			if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
				t.Errorf("could not unmarshal event. details: %v", err)
			}

			d, err := inputEvent.AWSLogs.Parse()
			if err != nil {
				if !test.expectError {
					t.Errorf("unexpected error: %+v", err)
				}

				if !reflect.DeepEqual(test.expectCloudwatchLogsData, d) {
					t.Errorf("expected: %+v, received: %v", test.expectCloudwatchLogsData, d)
				}
			}

			if err == nil && test.expectError {
				t.Errorf("expected error")
			}
		})
	}
}
