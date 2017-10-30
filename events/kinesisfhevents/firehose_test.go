// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package kinesisfhevents

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func testFirehoseEventMarshaling(t *testing.T) {
	testMarshaling(t, "./test-files/kinesis-firehose-event.json")
}

func testFirehoseResponseMarshaling(t *testing.T) {
	testMarshaling(t, "./test-files/kinesis-firehose-response.json")
}

func testMarshaling(t *testing.T, jsonFile string) {

	// 1. read JSON from file
	inputJson := readJsonFromFile(t, jsonFile)

	// 2. de-serialize into Go object
	var inputEvent KinesisFirehoseEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. serialize to JSON
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. check result
	test.AssertJsonsEqual(t, inputJson, outputJson)
}

func readJsonFromFile(t *testing.T, jsonFile string) []byte {
	inputJson, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJson
}

func TestSampleTransformation(t *testing.T) {

	inputJson := readJsonFromFile(t, "./test-files/kinesis-firehose-event.json")

	// de-serialize into Go object
	var inputEvent KinesisFirehoseEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	response := toUpperHandler(context.TODO(), inputEvent)

	inputString := string(inputEvent.Records[0].Data)
	expectedString := strings.ToUpper(inputString)
	actualString := string(response.Records[0].Data)
	assert.Equal(t, actualString, expectedString)
}

func toUpperHandler(ctx context.Context, evnt KinesisFirehoseEvent) KinesisFirehoseResponse {
	var response KinesisFirehoseResponse

	for _, record := range evnt.Records {
		// Transform data: ToUpper the data
		var transformedRecord FirehoseResponseRecord
		transformedRecord.RecordId = record.RecordId
		transformedRecord.Result = TransformedStateOk
		transformedRecord.Data = []byte(strings.ToUpper(string(record.Data)))

		response.Records = append(response.Records, transformedRecord)
	}

	return response
}

func TestKinesisFirehoseMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, KinesisFirehoseEvent{})
}
