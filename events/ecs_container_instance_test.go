// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestECSContainerInstanceEventMarshaling(t *testing.T) {
	// 1. read JSON from file
	inputJSON := test.ReadJSONFromFile(t, "./testdata/ecs-container-instance-state-change.json")

	// 2. de-serialize into Go object
	var inputEvent ECSContainerInstanceEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// 3. Verify values populated into Go Object, at least one validation per data type
	assert.Equal(t, "0", inputEvent.Version)
	assert.Equal(t, "8952ba83-7be2-4ab5-9c32-6687532d15a2", inputEvent.ID)
	assert.Equal(t, "ECS Container Instance State Change", inputEvent.DetailType)
	assert.Equal(t, "aws.ecs", inputEvent.Source)
	assert.Equal(t, "111122223333", inputEvent.Account)
	assert.Equal(t, "us-east-1", inputEvent.Region)
	assert.Equal(t, "arn:aws:ecs:us-east-1:111122223333:container-instance/b54a2a04-046f-4331-9d74-3f6d7f6ca315", inputEvent.Resources[0])
	testTime, err := time.Parse(time.RFC3339, "2016-12-06T16:41:06Z")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	assert.Equal(t, testTime, inputEvent.Time)

	var detail = inputEvent.Detail
	assert.True(t, detail.AgentConnected)
	assert.Equal(t, "com.amazonaws.ecs.capability.logging-driver.syslog", detail.Attributes[0].Name)
	assert.Equal(t, "arn:aws:ecs:us-east-1:111122223333:cluster/default", detail.ClusterARN)
	assert.Equal(t, "arn:aws:ecs:us-east-1:111122223333:container-instance/b54a2a04-046f-4331-9d74-3f6d7f6ca315", detail.ContainerInstanceARN)
	assert.Equal(t, "i-f3a8506b", detail.EC2InstanceID)
	assert.Equal(t, "CPU", detail.RegisteredResources[0].Name)
	assert.Equal(t, "INTEGER", detail.RegisteredResources[0].Type)
	assert.Equal(t, 2048, detail.RegisteredResources[0].IntegerValue)
	assert.Equal(t, "MEMORY", detail.RegisteredResources[1].Name)
	assert.Equal(t, "INTEGER", detail.RegisteredResources[1].Type)
	assert.Equal(t, 3767, detail.RegisteredResources[1].IntegerValue)
	assert.Equal(t, "PORTS", detail.RegisteredResources[2].Name)
	assert.Equal(t, "STRINGSET", detail.RegisteredResources[2].Type)
	assert.Equal(t, []*string{ptr("22"), ptr("2376"), ptr("2375"), ptr("51678"), ptr("51679")}, detail.RegisteredResources[2].StringSetValue)
	assert.Equal(t, "PORTS_UDP", detail.RegisteredResources[3].Name)
	assert.Equal(t, "STRINGSET", detail.RegisteredResources[3].Type)
	assert.Equal(t, []*string{}, detail.RegisteredResources[3].StringSetValue)
	assert.Equal(t, "CPU", detail.RemainingResources[0].Name)
	assert.Equal(t, "INTEGER", detail.RemainingResources[0].Type)
	assert.Equal(t, 1988, detail.RemainingResources[0].IntegerValue)
	assert.Equal(t, "MEMORY", detail.RemainingResources[1].Name)
	assert.Equal(t, "INTEGER", detail.RemainingResources[1].Type)
	assert.Equal(t, 767, detail.RemainingResources[1].IntegerValue)
	assert.Equal(t, "PORTS", detail.RemainingResources[2].Name)
	assert.Equal(t, "STRINGSET", detail.RemainingResources[2].Type)
	assert.Equal(t, []*string{ptr("22"), ptr("2376"), ptr("2375"), ptr("51678"), ptr("51679")}, detail.RemainingResources[2].StringSetValue)
	assert.Equal(t, "PORTS_UDP", detail.RemainingResources[3].Name)
	assert.Equal(t, "STRINGSET", detail.RemainingResources[3].Type)
	assert.Equal(t, []*string{}, detail.RemainingResources[3].StringSetValue)
	assert.Equal(t, "ACTIVE", detail.Status)
	assert.Equal(t, 14801, detail.Version)
	assert.Equal(t, "aebcbca", detail.VersionInfo.AgentHash)
	assert.Equal(t, "1.13.0", detail.VersionInfo.AgentVersion)
	assert.Equal(t, "DockerVersion: 1.11.2", detail.VersionInfo.DockerVersion)
	testUpdateTime, err := time.Parse(time.RFC3339, "2016-12-06T16:41:06Z")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	assert.Equal(t, testUpdateTime, inputEvent.Time)

	// 4. serialize to JSON
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	// 5. check result
	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func ptr(s string) *string {
	return &s
}

func TestECSContainerInstanceMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, ECSContainerInstanceEvent{})
}
