// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSecondsEpochTime(t *testing.T) {

	input := []byte("1480641523.476")

	var epoch SecondsEpochTime
	if err := json.Unmarshal(input, &epoch); err != nil {
		t.Errorf("unmarshal failed. details: %v", err)
	}

	utc := epoch.UTC()
	assert.Equal(t, 2016, utc.Year())
	assert.Equal(t, time.Month(12), utc.Month())
	assert.Equal(t, 2, utc.Day())
	assert.Equal(t, 1, utc.Hour())
	assert.Equal(t, 18, utc.Minute())
	assert.Equal(t, 43, utc.Second())
	assert.Equal(t, 476, utc.Nanosecond()/1000000)
}

func TestMarshalSecondsEpochTime(t *testing.T) {
	timestamp := SecondsEpochTime{time.Date(2016, time.Month(12), 2, 1, 18, 43, 476*1000000, time.UTC)}

	marshaled, err := json.Marshal(timestamp)
	if err != nil {
		t.Errorf("marshal failed. details: %v", err)
	}

	assert.Equal(t, "1480641523.476", string(marshaled))
}

func TestUnmarshalMilliSecondsEpochTime(t *testing.T) {

	input := []byte("1480641523476")

	var epoch MilliSecondsEpochTime
	if err := json.Unmarshal(input, &epoch); err != nil {
		t.Errorf("unmarshal failed. details: %v", err)
	}

	utc := epoch.UTC()
	assert.Equal(t, 2016, utc.Year())
	assert.Equal(t, time.Month(12), utc.Month())
	assert.Equal(t, 2, utc.Day())
	assert.Equal(t, 1, utc.Hour())
	assert.Equal(t, 18, utc.Minute())
	assert.Equal(t, 43, utc.Second())
	assert.Equal(t, 476, utc.Nanosecond()/1000000)
}

func TestMarshalMilliSecondsEpochTime(t *testing.T) {
	timestamp := MilliSecondsEpochTime{time.Date(2016, time.Month(12), 2, 1, 18, 43, 476*1000000, time.UTC)}

	marshaled, err := json.Marshal(timestamp)
	if err != nil {
		t.Errorf("marshal failed. details: %v", err)
	}

	assert.Equal(t, "1480641523476", string(marshaled))
}
