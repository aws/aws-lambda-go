// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestS3EventMarshaling(t *testing.T) {

	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/s3-event.json")

	// 2. de-serialize into Go object
	var inputEvent S3Event
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 4. read expected output JSON from file
	exepectedOutputJSON := test.ReadJSONFromFile(t, "./testdata/s3-event-with-decoded.json")

	// 5. check result
	assert.JSONEq(t, string(exepectedOutputJSON), string(outputJSON))
}

func TestS3TestEventMarshaling(t *testing.T) {
	inputJSON := []byte(`{
	    "Service" :"Amazon S3",
	    "Event": "s3:TestEvent",
	    "Time": "2019-02-04T19:34:46.985Z",
	    "Bucket": "bmoffatt",
	    "RequestId": "7BA1940DC6AF888B",
	    "HostId": "q1YDbiaMjllP0m+Lcy6cKKgxNrMLFJ9zCrZUFBqHGTG++0nXvnTDIGC7q2/QPAsJg86E8gI7y9U="
	}`)
	var inputEvent S3TestEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestS3MarshalingMalformedJSON(t *testing.T) {
	test.TestMalformedJson(t, S3Event{})
}

func TestS3GlacierEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/s3-glacier-event.json")

	// 2. de-serialize into Go object
	var inputEvent S3Event
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. verify glacierEventData is correctly parsed
	if inputEvent.Records[0].GlacierEventData == nil {
		t.Error("glacierEventData should not be nil for glacier restore events")
	}

	// 4. verify restoreEventData is correctly parsed
	if inputEvent.Records[0].GlacierEventData.RestoreEventData == nil {
		t.Error("restoreEventData should not be nil")
	}

	// 5. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 6. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestS3RestoreEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/s3-restore-event.json")

	// 2. de-serialize into Go object
	var inputEvent S3Event
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. verify restoreEventData is correctly parsed
	if inputEvent.Records[0].RestoreEventData == nil {
		t.Error("restoreEventData should not be nil")
	}

	// 4. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 5. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestS3IntelligentTieringEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/s3-intelligenttier-event.json")

	// 2. de-serialize into Go object
	var inputEvent S3Event
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. verify intelligentTieringEventData is correctly parsed
	if inputEvent.Records[0].IntelligentTieringEventData == nil {
		t.Error("intelligentTieringEventData should not be nil for intelligent tiering events")
	}

	// 4. verify destinationAccessTier is correctly parsed
	if inputEvent.Records[0].IntelligentTieringEventData.DestinationAccessTier == "" {
		t.Error("destinationAccessTier should not be empty")
	}

	// 5. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 6. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestS3LifecycleEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/s3-lifecycle-event.json")

	// 2. de-serialize into Go object
	var inputEvent S3Event
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. verify lifecycleEventData is correctly parsed
	if inputEvent.Records[0].LifecycleEventData == nil {
		t.Error("lifecycleEventData should not be nil for lifecycle events")
	}

	// 4. verify transitionEventData is correctly parsed
	if inputEvent.Records[0].LifecycleEventData.TransitionEventData == nil {
		t.Error("transitionEventData should not be nil")
	}

	// 5. verify destinationStorageClass is correctly parsed
	if inputEvent.Records[0].LifecycleEventData.TransitionEventData.DestinationStorageClass == "" {
		t.Error("destinationStorageClass should not be empty")
	}

	// 6. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 7. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestS3ReplicationEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/s3-replication-event.json")

	// 2. de-serialize into Go object
	var inputEvent S3Event
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. verify replicationEventData is correctly parsed
	if inputEvent.Records[0].ReplicationEventData == nil {
		t.Error("replicationEventData should not be nil for replication events")
	}

	// 4. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 5. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}
