package test

import (
	"io/ioutil"
	"testing"
)

// ReadJSONFromFile reads a given input file to JSON
func ReadJSONFromFile(t *testing.T, inputFile string) []byte {
	inputJSON, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJSON
}
