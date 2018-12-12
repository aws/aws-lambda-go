package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestCloudFrontRequestMarshaling(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("./testdata/cloudfront-request.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputEvent CloudFrontRequest
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}
	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestCloudFrontRequestMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CloudFrontRequest{})
}

func TestCloudFrontResponseMarshaling(t *testing.T) {
	inputJSON, err := ioutil.ReadFile("./testdata/cloudfront-response.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var inputEvent CloudFrontResponse
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestCloudFrontResponseMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CloudFrontResponse{})
}
