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
