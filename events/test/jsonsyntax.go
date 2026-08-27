// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package test

import (
	"encoding/json"
	"errors"
	"testing"
)

// unmarshalJSON is the function used by TestMalformedJson to deserialize JSON.
// It defaults to json.Unmarshal and can be overridden in tests within this package.
var unmarshalJSON = json.Unmarshal

// nolint: staticcheck
func TestMalformedJson(t testing.TB, objectToDeserialize interface{}) {
	// 1. read JSON from file
	inputJson := GetMalformedJson()

	// 2. de-serialize into Go object
	err := unmarshalJSON(inputJson, objectToDeserialize)
	if err == nil {
		t.Errorf("unmarshal should have failed but succeeded instead")
	}

	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("unmarshal should have returned a json.SyntaxError")
	}
}
