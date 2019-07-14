package events

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestUnmarshalCodeDeployEvent(t *testing.T) {
	tests := []struct {
		input  string
		expect CodeDeployEvent
	}{
		{
			input: "testdata/codedeploy-deployment-event.json",
			expect: CodeDeployEvent{
				AccountID:  "123456789012",
				Region:     "us-east-1",
				DetailType: CodeDeployDeploymentEventDetailType,
				Source:     CodeDeployEventSource,
				Version:    "0",
				Time:       time.Date(2016, 6, 30, 22, 6, 31, 0, time.UTC),
				ID:         "c071bfbf-83c4-49ca-a6ff-3df053957145",
				Resources: []string{
					"arn:aws:codedeploy:us-east-1:123456789012:application:myApplication",
					"arn:aws:codedeploy:us-east-1:123456789012:deploymentgroup:myApplication/myDeploymentGroup",
				},
				Detail: CodeDeployEventDetail{
					InstanceGroupID: "9fd2fbef-2157-40d8-91e7-6845af69e2d2",
					InstanceID:      "",
					Region:          "us-east-1",
					Application:     "myApplication",
					DeploymentID:    "d-123456789",
					State:           CodeDeployDeploymentStateSuccess,
					DeploymentGroup: "myDeploymentGroup",
				},
			},
		},
		{
			input: "testdata/codedeploy-instance-event.json",
			expect: CodeDeployEvent{
				AccountID:  "123456789012",
				Region:     "us-east-1",
				DetailType: CodeDeployInstanceEventDetailType,
				Source:     CodeDeployEventSource,
				Version:    "0",
				Time:       time.Date(2016, 6, 30, 23, 18, 50, 0, time.UTC),
				ID:         "fb1d3015-c091-4bf9-95e2-d98521ab2ecb",
				Resources: []string{
					"arn:aws:ec2:us-east-1:123456789012:instance/i-0000000aaaaaaaaaa",
					"arn:aws:codedeploy:us-east-1:123456789012:deploymentgroup:myApplication/myDeploymentGroup",
					"arn:aws:codedeploy:us-east-1:123456789012:application:myApplication",
				},
				Detail: CodeDeployEventDetail{
					InstanceGroupID: "8cd3bfa8-9e72-4cbe-a1e5-da4efc7efd49",
					InstanceID:      "i-0000000aaaaaaaaaa",
					Region:          "us-east-1",
					Application:     "myApplication",
					DeploymentID:    "d-123456789",
					State:           CodeDeployDeploymentStateSuccess,
					DeploymentGroup: "myDeploymentGroup",
				},
			},
		},
	}

	for _, testcase := range tests {
		data, err := ioutil.ReadFile(testcase.input)
		require.NoError(t, err)

		var actual CodeDeployEvent
		require.NoError(t, json.Unmarshal(data, &actual))

		require.Equal(t, testcase.expect, actual)
	}
}
