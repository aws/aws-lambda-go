// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestECRImageActionEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/ecr-image-push-event.json")

	// 2. de-serialize into Go object
	var inputEvent ECRImageActionEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. Verify values populated into Go Object, at least one validation per data type
	assert.Equal(t, "0", inputEvent.Version)
	assert.Equal(t, "13cde686-328b-6117-af20-0e5566167482", inputEvent.ID)
	assert.Equal(t, "ECR Image Action", inputEvent.DetailType)
	assert.Equal(t, "aws.ecr", inputEvent.Source)
	assert.Equal(t, "123456789012", inputEvent.Account)
	assert.Equal(t, "us-west-2", inputEvent.Region)
	assert.Empty(t, inputEvent.Resources)

	testTime, err := time.Parse(time.RFC3339, "2019-11-16T01:54:34Z")
	assert.Equal(t, testTime, inputEvent.Time)

	var detail = inputEvent.Detail
	assert.Equal(t, "SUCCESS", detail.Result)
	assert.Equal(t, "my-repository-name", detail.RepositoryName)
	assert.Equal(t, "sha256:7f5b2640fe6fb4f46592dfd3410c4a79dac4f89e4782432e0378abcd1234", detail.ImageDigest)
	assert.Equal(t, "latest", detail.ImageTag)

	// 4. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 5. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestECRPushMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ECRImageActionEvent{})
}
