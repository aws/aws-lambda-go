// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package test

import (
	"encoding/json"
	"testing"
)

func TestMalformedJson(t *testing.T, objectToDeserialize interface{}) {
	// 1. read JSON from file
	inputJson := GetMalformedJson()

	// 2. de-serialize into Go object
	err := json.Unmarshal(inputJson, objectToDeserialize)
	if err == nil {
		t.Errorf("unmarshal should have failed but succeeded instead")
	}

	_, isSyntaxError := err.(*json.SyntaxError)
	if !isSyntaxError {
		t.Errorf("unmarshal should have returned a json.SyntaxError")
	}
}
