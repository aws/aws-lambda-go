// Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAMPolicyStatementConditionRoundTrip(t *testing.T) {
	// A real IAM policy statement (as returned by an API Gateway custom
	// authorizer) can carry a Condition block; it must survive a round trip.
	const input = `{"Action":["execute-api:Invoke"],"Effect":"Deny","Resource":["arn:aws:execute-api:us-east-1:123456789012:apiId/stage/GET/path"],"Condition":{"IpAddress":{"aws:SourceIp":["1.2.3.4/24"]}}}`

	var stmt IAMPolicyStatement
	require.NoError(t, json.Unmarshal([]byte(input), &stmt))

	output, err := json.Marshal(stmt)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(output))

	// A JSON array condition value unmarshals to []interface{}.
	assert.Equal(t, []interface{}{"1.2.3.4/24"}, stmt.Condition["IpAddress"]["aws:SourceIp"])
}

func TestIAMPolicyStatementScalarCondition(t *testing.T) {
	// A condition value may also be a single scalar rather than an array, which
	// is why the value type has to be interface{} and not []string.
	const input = `{"Action":["execute-api:Invoke"],"Effect":"Allow","Resource":["*"],"Condition":{"Bool":{"aws:MultiFactorAuthPresent":"true"}}}`

	var stmt IAMPolicyStatement
	require.NoError(t, json.Unmarshal([]byte(input), &stmt))
	assert.Equal(t, "true", stmt.Condition["Bool"]["aws:MultiFactorAuthPresent"])

	output, err := json.Marshal(stmt)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(output))
}

func TestIAMPolicyStatementWithoutConditionOmitted(t *testing.T) {
	// A statement without a condition must not serialize "Condition":null,
	// otherwise the custom-authorizer response golden round trip breaks.
	output, err := json.Marshal(IAMPolicyStatement{
		Action:   []string{"execute-api:Invoke"},
		Effect:   "Allow",
		Resource: []string{"*"},
	})
	require.NoError(t, err)
	assert.JSONEq(t, `{"Action":["execute-api:Invoke"],"Effect":"Allow","Resource":["*"]}`, string(output))
}
