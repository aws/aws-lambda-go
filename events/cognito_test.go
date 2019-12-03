// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestCognitoEventMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestCognitoMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEvent{})
}

func TestCognitoEventUserPoolsPreSignupMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event-userpools-presignup.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEventUserPoolsPreSignup
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestCognitoUserPoolsPreSignupMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsPreSignup{})
}

func TestCognitoEventUserPoolsPreAuthenticationMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event-userpools-preauthentication.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEventUserPoolsPreAuthentication
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestCognitoUserPoolsPreAuthenticationMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsPreAuthentication{})
}

func TestCognitoEventUserPoolsPostConfirmationMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event-userpools-postconfirmation.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEventUserPoolsPostConfirmation
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestCognitoEventUserPoolsPreTokenGenMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsPreTokenGen{})
}

func TestCognitoEventUserPoolsPreTokenGenMarshaling(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event-userpools-pretokengen.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEventUserPoolsPreTokenGen
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestCognitoEventUserPoolsDefineAuthChallengeMarshaling(t *testing.T) {
	var inputEvent CognitoEventUserPoolsDefineAuthChallenge
	test.AssertJsonFile(t, "./testdata/cognito-event-userpools-define-auth-challenge.json", &inputEvent)
}

func TestCognitoEventUserPoolsDefineAuthChallengeMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsDefineAuthChallenge{})
}

func TestCognitoEventUserPoolsCreateAuthChallengeMarshaling(t *testing.T) {
	var inputEvent CognitoEventUserPoolsCreateAuthChallenge
	test.AssertJsonFile(t, "./testdata/cognito-event-userpools-create-auth-challenge.json", &inputEvent)
}

func TestCognitoEventUserPoolsCreateAuthChallengeMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsCreateAuthChallenge{})
}

func TestCognitoEventUserPoolsVerifyAuthChallengeMarshaling(t *testing.T) {
	var inputEvent CognitoEventUserPoolsVerifyAuthChallenge
	test.AssertJsonFile(t, "./testdata/cognito-event-userpools-verify-auth-challenge.json", &inputEvent)
}

func TestCognitoEventUserPoolsVerifyAuthChallengeMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsVerifyAuthChallenge{})
}

func TestCognitoEventUserPoolsPostAuthenticationMarshaling(t *testing.T) {

	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event-userpools-postauthentication.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEventUserPoolsPostAuthentication
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestCognitoEventUserPoolsMigrateUserMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsMigrateUser{})
}

func TestCognitoEventUserPoolsMigrateUserMarshaling(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event-userpools-migrateuser.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEventUserPoolsMigrateUser
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}
	test.AssertJsonsEqual(t, inputJSON, outputJSON)
}

func TestCognitoEventUserPoolsCustomMessageMarshaling(t *testing.T) {
	// read json from file
	inputJSON, err := ioutil.ReadFile("./testdata/cognito-event-userpools-custommessage.json")
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	// de-serialize into CognitoEvent
	var inputEvent CognitoEventUserPoolsCustomMessage
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	// serialize to json
	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestCognitoUserPoolsCustomMessageMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, CognitoEventUserPoolsCustomMessage{})
}
