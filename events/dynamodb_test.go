// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestDynamoDBEventMarshaling(t *testing.T) {

	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/dynamodb-event.json")

	// 2. de-serialize into Go object
	var inputEvent DynamoDBEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestDynamoDBEventMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, &DynamoDBEvent{})
}

func TestDynamoDBStreamRecordCreationDateTimePrecision(t *testing.T) {
	testCases := []struct {
		name      string
		epoch     int64
		precision string
		want      time.Time
	}{
		{"millisecond", 1731101300058, "MILLISECOND", time.Date(2024, time.November, 8, 21, 28, 20, 58000000, time.UTC)},
		{"microsecond", 1731101300058336, "MICROSECOND", time.Date(2024, time.November, 8, 21, 28, 20, 58336000, time.UTC)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inputJSON := []byte(fmt.Sprintf(`{
				"ApproximateCreationDateTime": %d,
				"ApproximateCreationDateTimePrecision": %q,
				"SequenceNumber": "1",
				"SizeBytes": 1,
				"StreamViewType": "NEW_IMAGE"
			}`, testCase.epoch, testCase.precision))

			var inputRecord DynamoDBStreamRecord
			if err := json.Unmarshal(inputJSON, &inputRecord); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.want, inputRecord.ApproximateCreationDateTime.UTC())
			assert.Equal(t, testCase.precision, inputRecord.ApproximateCreationDateTimePrecision)

			outputJSON, err := json.Marshal(inputRecord)
			if err != nil {
				t.Fatal(err)
			}
			assert.JSONEq(t, string(inputJSON), string(outputJSON))
		})
	}
}

func TestDynamoDBStreamRecordCreationDateTimeRange(t *testing.T) {
	testCases := []struct {
		name      string
		precision string
		epoch     int64
		want      time.Time
	}{
		{"millisecond min", "MILLISECOND", -9223372036854775807 - 1, time.UnixMilli(-9223372036854775807 - 1)},
		{"millisecond max", "MILLISECOND", 9223372036854775807, time.UnixMilli(9223372036854775807)},
		{"microsecond min", "MICROSECOND", -9223372036854775807 - 1, time.UnixMicro(-9223372036854775807 - 1)},
		{"microsecond max", "MICROSECOND", 9223372036854775807, time.UnixMicro(9223372036854775807)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inputJSON := []byte(fmt.Sprintf(`{
				"ApproximateCreationDateTime": %d,
				"ApproximateCreationDateTimePrecision": %q
			}`, testCase.epoch, testCase.precision))

			var inputRecord DynamoDBStreamRecord
			if err := json.Unmarshal(inputJSON, &inputRecord); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.want, inputRecord.ApproximateCreationDateTime.Time)

			outputJSON, err := json.Marshal(inputRecord)
			if err != nil {
				t.Fatal(err)
			}
			var outputRecord map[string]json.RawMessage
			if err := json.Unmarshal(outputJSON, &outputRecord); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, strconv.FormatInt(testCase.epoch, 10), string(outputRecord["ApproximateCreationDateTime"]))
			assert.Equal(t, 1, strings.Count(string(outputJSON), `"ApproximateCreationDateTime":`))
		})
	}
}

func TestDynamoDBStreamRecordRejectsUnknownCreationDateTimePrecisionOnUnmarshal(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{
			"ApproximateCreationDateTime": 1731101300058336,
			"ApproximateCreationDateTimePrecision": "NANOSECOND"
		}`),
		[]byte(`{
			"ApproximateCreationDateTimePrecision": "NANOSECOND"
		}`),
	}

	for _, inputJSON := range testCases {
		var inputRecord DynamoDBStreamRecord
		if err := json.Unmarshal(inputJSON, &inputRecord); err == nil {
			t.Fatal("expected unsupported precision error")
		}
	}
}

func TestDynamoDBStreamRecordRejectsUnknownCreationDateTimePrecisionOnMarshal(t *testing.T) {
	record := DynamoDBStreamRecord{
		ApproximateCreationDateTime:          SecondsEpochTime{time.Unix(0, 0)},
		ApproximateCreationDateTimePrecision: "NANOSECOND",
	}

	if _, err := json.Marshal(record); err == nil {
		t.Fatal("expected unsupported precision error")
	}
}

func TestDynamoDBStreamRecordRejectsInvalidPreciseCreationDateTime(t *testing.T) {
	testCases := []struct {
		name      string
		precision string
		epoch     string
	}{
		{"fractional milliseconds", "MILLISECOND", "1.5"},
		{"fractional microseconds", "MICROSECOND", "1.5"},
		{"millisecond overflow", "MILLISECOND", "9223372036854775808"},
		{"microsecond underflow", "MICROSECOND", "-9223372036854775809"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inputJSON := []byte(fmt.Sprintf(`{
				"ApproximateCreationDateTime": %s,
				"ApproximateCreationDateTimePrecision": %q
			}`, testCase.epoch, testCase.precision))

			var inputRecord DynamoDBStreamRecord
			if err := json.Unmarshal(inputJSON, &inputRecord); err == nil {
				t.Fatal("expected invalid timestamp error")
			}
		})
	}
}

func TestDynamoDBTimeWindowEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/dynamodb-time-window-event.json")

	// 2. de-serialize into Go object
	var inputEvent DynamoDBTimeWindowEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestDynamoDBTimeWindowEventMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, &DynamoDBTimeWindowEvent{})
}
