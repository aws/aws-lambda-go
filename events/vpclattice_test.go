package events

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVPCLatticeRequestV1Marshalling(t *testing.T) {
	test.AssertJsonFile(t, "./testdata/vpclattice-v1-request.json", &VPCLatticeRequestV1{})
}

func TestVPCLatticeRequestV1MalformedJson(t *testing.T) {
	test.TestMalformedJson(t, &VPCLatticeRequestV1{})
}

func TestVPCLatticeRequestV2Marshalling(t *testing.T) {
	test.AssertJsonFile(t, "./testdata/vpclattice-v2-request.json", &VPCLatticeRequestV2{})
}

func TestVPCLatticeRequestV2MalformedJson(t *testing.T) {
	test.TestMalformedJson(t, &VPCLatticeRequestV2{})
}

func TestVPCLatticeRequestV1IncludesFalseBase64Flag(t *testing.T) {
	var request VPCLatticeRequestV1
	require.NoError(t, json.Unmarshal([]byte(`{"is_base64_encoded":false}`), &request))
	output, err := json.Marshal(request)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(output, &fields))
	assert.Equal(t, json.RawMessage("false"), fields["is_base64_encoded"])
}

func TestVPCLatticeRequestV2OptionalIdentityFields(t *testing.T) {
	const input = `{"version":"2.0","path":"/","method":"GET","headers":{},"body":"","requestContext":{"serviceNetworkArn":"network","serviceArn":"service","targetGroupArn":"target","identity":{"type":"AWS_IAM"},"region":"us-east-1","timeEpoch":"1695799509392227"}}`
	var request VPCLatticeRequestV2
	require.NoError(t, json.Unmarshal([]byte(input), &request))
	output, err := json.Marshal(request)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(output))
}

func TestVPCLatticeResponse(t *testing.T) {
	test.AssertJsonFile(t, "./testdata/vpclattice-response.json", &VPCLatticeResponse{})
}

func TestVPCLatticeResponseWithoutBody(t *testing.T) {
	response := VPCLatticeResponse{StatusCode: 204, Headers: map[string]string{}}
	output, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, `{"isBase64Encoded":false,"statusCode":204,"headers":{}}`, string(output))
}
