// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestECRScanEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJson := test.ReadJSONFromFile(t, "./testdata/ecr-image-scan-event.json")

	// 2. de-serialize into Go object
	var inputEvent ECRScanEvent
	if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. Verify values populated into Go Object, at least one validation per data type
	assert.Equal(t, "0", inputEvent.Version)
	assert.Equal(t, "01234567-0123-0123-0123-012345678901", inputEvent.ID)
	assert.Equal(t, "ECR Image Scan", inputEvent.DetailType)
	assert.Equal(t, "aws.ecr", inputEvent.Source)
	assert.Equal(t, "123456789012", inputEvent.Account)
	assert.Equal(t, "2019-10-30T21:32:27Z", inputEvent.Time)
	assert.Equal(t, "eu-north-1", inputEvent.Region)
	assert.Equal(t, "arn:aws:ecr:eu-north-1:123456789012:repository/tribble-image-scan-test", inputEvent.Resources[0])

	var detail = inputEvent.Detail
	assert.Equal(t, "COMPLETE", detail.ScanStatus)
	assert.Equal(t, "tribble-image-scan-test", detail.RepositoryName)
	assert.Equal(t, "sha256:d4a96ee9443e641fc100e763a0c10928720b50c6e3ea3342d05d7c3435fc5355", detail.ImageDigest)
	assert.Equal(t, "1572471135", detail.ImageTags[0])
	assert.Equal(t, int64(10), detail.FindingSeverityCounts.Critical)
	assert.Equal(t, int64(2), detail.FindingSeverityCounts.High)
	assert.Equal(t, int64(9), detail.FindingSeverityCounts.Medium)
	assert.Equal(t, int64(3), detail.FindingSeverityCounts.Low)
	assert.Equal(t, int64(0), detail.FindingSeverityCounts.Informational)
	assert.Equal(t, int64(0), detail.FindingSeverityCounts.Undefined)

	// 4. serialize to JSON
	outputJson, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 5. check result
	assert.JSONEq(t, string(inputJson), string(outputJson))
}

func TestECRScanMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ECRScanEvent{})
}
