// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Asserts two JSON files are semantically equal
// (ignores white-space and attribute order)
func AssertJsonsEqual(t *testing.T, expectedJson []byte, actualJson []byte) {
	result, err := areJsonsEqual(expectedJson, actualJson)
	if err != nil {
		t.Errorf("json comparison failed. EXPECTED json: %s\n\n ACTUAL json: %s\n\n . Error details: %s",
			string(expectedJson), string(actualJson), err)
	}

	assert.True(t, result)
}

func areJsonsEqual(expectedJson []byte, actualJson []byte) (bool, error) {

	expectedMap, err := unmarshalToMap(expectedJson)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal to map. details: %s", err)
	}

	actualMap, err := unmarshalToMap(actualJson)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal to map. details: %s", err)
	}

	return areMapsEqual(expectedMap, actualMap), nil
}

func unmarshalToMap(inputJson []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(inputJson, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func areMapsEqual(map1 map[string]interface{}, map2 map[string]interface{}) bool {
	if map1 == nil || map2 == nil {
		return false
	}

	if len(map1) != len(map2) {
		return false
	}

	for k, v1 := range map1 {
		v2, found := map2[k]
		if !found {
			return false
		}

		switch v1 := v1.(type) {
		case []interface{}:
			// compare lists
			if !areListsEqual(v1, v2.([]interface{})) {
				return false
			}
		case map[string]interface{}:
			// compare maps
			if !areMapsEqual(v1, v2.(map[string]interface{})) {
				return false
			}
		default:
			if v1 != v2 {
				return false
			}
		}
	}

	return true
}

func areListsEqual(list1 []interface{}, list2 []interface{}) bool {
	if list1 == nil || list2 == nil {
		return false
	}

	if len(list1) != len(list2) {
		return false
	}

	for index := range list1 {

		v1 := list1[index]
		v2 := list2[index]

		switch v1 := v1.(type) {
		case []interface{}:
			// compare lists
			if !areListsEqual(v1, v2.([]interface{})) {
				return false
			}
		case map[string]interface{}:
			// compare maps
			if !areMapsEqual(v1, v2.(map[string]interface{})) {
				return false
			}
		default:
			if v1 != v2 {
				return false
			}
		}
	}

	return true
}
