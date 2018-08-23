package events

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCloudwatchLogs(t *testing.T) {
	for _, test := range []struct {
		name                      string
		eventJson                 string
		expectError               bool
		expectCloudwatchEventData CloudwatchLogsEvent
	}{
		{"Well formed cloudwatch event",
			"./testdata/cloudwatch-logs-event.json",
			false,
			CloudwatchLogsEvent{
				AWSLogs: CloudwatchLogsRawData{
					Data: "H4sIAAAAAAAAAHWPwQqCQBCGX0Xm7EFtK+smZBEUgXoLCdMhFtKV3akI8d0bLYmibvPPN3wz00CJxmQnTO41whwWQRIctmEcB6sQbFC3CjW3XW8kxpOpP+OC22d1Wml1qZkQGtoMsScxaczKN3plG8zlaHIta5KqWsozoTYw3/djzwhpLwivWFGHGpAFe7DL68JlBUk+l7KSN7tCOEJ4M3/qOI49vMHj+zCKdlFqLaU2ZHV2a4Ct/an0/ivdX8oYc1UVX860fQDQiMdxRQEAAA==",
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			inputJson := readJsonFromFile(t, test.eventJson)

			var inputEvent CloudwatchLogsEvent
			err := json.Unmarshal(inputJson, &inputEvent)

			if err != nil && !test.expectError {
				t.Errorf("could not unmarshal event. details: %v", err)
			}

			if err == nil && test.expectError {
				t.Errorf("expected parse error")
			}

			if !reflect.DeepEqual(test.expectCloudwatchEventData, inputEvent) {
				t.Errorf("expected: %+v, received: %v", test.expectCloudwatchEventData, inputEvent)
			}
		})
	}
}

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
				LogEvents: []CloudwatchLogsLogEvent{
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
