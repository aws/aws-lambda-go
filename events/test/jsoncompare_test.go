// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertJsonCompareResult(t *testing.T, s1 string, s2 string, expected bool) {
	result, err := areJsonsEqual([]byte(s1), []byte(s2))
	if err != nil {
		t.Errorf("json comparison failed. details: %s", err)
	}

	assert.Equal(t, expected, result)
}

func assertJsonCompareFails(t *testing.T, s1 string, s2 string) {
	_, err := areJsonsEqual([]byte(s1), []byte(s2))
	if err == nil {
		t.Errorf("json comparison should have failed but succeeded instead")
	}
}

func TestJsonCompareReturnsFalseForMapsWithDifferentLengths(t *testing.T) {

	// maps
	assertJsonCompareResult(t, `{ "a": 1 }`, `{ "a": 1, "b": 2}`, false)
	assertJsonCompareResult(t, `{ "a": 1, "b": 2 }`, `{ "a": 1, "c": 2}`, false)
	assertJsonCompareResult(t, `{ "a": { "b" : 1 } }`, `{ "a": { "b" : 2 } }`, false)
	assertJsonCompareResult(t, `{ "a": 1, "b": 2 }`, `{ "b": 2, "a": 1}`, true)
	assertJsonCompareResult(t, `{ "a": { "b" : 1 } }`, `{ "a": { "b" : 1 } }`, true)

	// lists
	assertJsonCompareResult(t, `{ "a": [ 1 ] }`, `{ "a": [ 1, 2 ]}`, false)
	assertJsonCompareResult(t, `{ "a": [ 1, 2 ] }`, `{ "a": [ 1 ]}`, false)
	assertJsonCompareResult(t, `{ "a": [ 1, 2 ] }`, `{ "a": [ 2, 1 ]}`, false)
	assertJsonCompareResult(t, `{ "a": [ 1, 2 ] }`, `{ "a": [ 1, 2 ]}`, true)
	assertJsonCompareResult(t, `{ "a": [  ] }`, `{ "a": [ ]}`, true)

	// malformed inputs
	assertJsonCompareFails(t, `{ "a": 1}`, `{ "a: 1`)
	assertJsonCompareFails(t, `{ "a}`, `{ "a": 1`)
}
