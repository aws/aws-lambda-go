package events

import (
	"encoding/json"
	"reflect"
	"testing"

	tst "github.com/aws/aws-lambda-go/events/test"
)

func TestCloudwatchLogs(t *testing.T) {
	for _, test := range []struct {
		name                  string
		eventJSON             string
		expectError           bool
		expectEventBridgeData EventBridgeEventLogs
	}{
		{"Well formed cloudwatch event",
			"./testdata/cloudwatch-logs-event.json",
			false,
			EventBridgeEventLogs{
				AWSLogs: EventBridgeLogsRawData{
					Data: "H4sIAAAAAAAAAHWPwQqCQBCGX0Xm7EFtK+smZBEUgXoLCdMhFtKV3akI8d0bLYmibvPPN3wz00CJxmQnTO41whwWQRIctmEcB6sQbFC3CjW3XW8kxpOpP+OC22d1Wml1qZkQGtoMsScxaczKN3plG8zlaHIta5KqWsozoTYw3/djzwhpLwivWFGHGpAFe7DL68JlBUk+l7KSN7tCOEJ4M3/qOI49vMHj+zCKdlFqLaU2ZHV2a4Ct/an0/ivdX8oYc1UVX860fQDQiMdxRQEAAA==",
				},
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			inputJSON := tst.ReadJSONFromFile(t, test.eventJSON)

			var inputEvent EventBridgeEventLogs
			err := json.Unmarshal(inputJSON, &inputEvent)

			if err != nil && !test.expectError {
				t.Errorf("could not unmarshal event. details: %v", err)
			}

			if err == nil && test.expectError {
				t.Errorf("expected parse error")
			}

			if !reflect.DeepEqual(test.expectEventBridgeData, inputEvent) {
				t.Errorf("expected: %+v, received: %v", test.expectEventBridgeData, inputEvent)
			}
		})
	}
}

func TestCloudwatchLogsParse(t *testing.T) {
	for _, test := range []struct {
		name                     string
		eventJSON                string
		expectError              bool
		expectCloudwatchLogsData EventBridgeLogsData
	}{
		{"Well formed cloudwatch event",
			"./testdata/cloudwatch-logs-event.json",
			false,
			EventBridgeLogsData{
				Owner:     "123456789123",
				LogGroup:  "testLogGroup",
				LogStream: "testLogStream",
				SubscriptionFilters: []string{
					"testFilter",
				},
				MessageType: "DATA_MESSAGE",
				LogEvents: []EventBridgeLogEvent{
					{ID: "eventId1", Timestamp: 1440442987000, Message: "[ERROR] First test message"},
					{ID: "eventId2", Timestamp: 1440442987001, Message: "[ERROR], Second test message"},
				},
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			inputJSON := tst.ReadJSONFromFile(t, test.eventJSON)

			var inputEvent EventBridgeEventLogs
			if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
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
