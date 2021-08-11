package events

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestUnmarshalCodePipelineEvent(t *testing.T) {
	tests := []struct {
		input  string
		expect CodePipelineCloudWatchEvent
	}{
		{
			input: "testdata/codepipeline-action-execution-stage-change-event.json",
			expect: CodePipelineCloudWatchEvent{
				Version:    "0",
				ID:         "CWE-event-id",
				DetailType: "CodePipeline Action Execution State Change",
				Source:     "aws.codepipeline",
				AccountID:  "123456789012",
				Time:       time.Date(2017, 04, 22, 3, 31, 47, 0, time.UTC),
				Region:     "us-east-1",
				Resources: []string{
					"arn:aws:codepipeline:us-east-1:123456789012:pipeline:myPipeline",
				},
				Detail: CodePipelineEventDetail{
					Pipeline:    "myPipeline",
					Version:     1,
					ExecutionID: "01234567-0123-0123-0123-012345678901",
					Stage:       "Prod",
					Action:      "myAction",
					State:       "STARTED",
					Region:      "us-west-2",
					Type: CodePipelineEventDetailType{
						Owner:    "AWS",
						Category: "Deploy",
						Provider: "CodeDeploy",
						Version:  "1",
					},
				},
			},
		},
		{
			input: "testdata/codepipeline-execution-stage-change-event.json",
			expect: CodePipelineCloudWatchEvent{
				Version:    "0",
				ID:         "CWE-event-id",
				DetailType: "CodePipeline Stage Execution State Change",
				Source:     "aws.codepipeline",
				AccountID:  "123456789012",
				Time:       time.Date(2017, 04, 22, 3, 31, 47, 0, time.UTC),
				Region:     "us-east-1",
				Resources: []string{
					"arn:aws:codepipeline:us-east-1:123456789012:pipeline:myPipeline",
				},
				Detail: CodePipelineEventDetail{
					Pipeline:    "myPipeline",
					Version:     1,
					ExecutionID: "01234567-0123-0123-0123-012345678901",
					State:       "STARTED",
				},
			},
		},
		{
			input: "testdata/codepipeline-execution-state-change-event.json",
			expect: CodePipelineCloudWatchEvent{
				Version:    "0",
				ID:         "CWE-event-id",
				DetailType: "CodePipeline Pipeline Execution State Change",
				Source:     "aws.codepipeline",
				AccountID:  "123456789012",
				Time:       time.Date(2017, 04, 22, 3, 31, 47, 0, time.UTC),
				Region:     "us-east-1",
				Resources: []string{
					"arn:aws:codepipeline:us-east-1:123456789012:pipeline:myPipeline",
				},
				Detail: CodePipelineEventDetail{
					Pipeline:    "myPipeline",
					Version:     1,
					ExecutionID: "01234567-0123-0123-0123-012345678901",
					State:       "STARTED",
				},
			},
		},
	}

	for _, testcase := range tests {
		data, err := ioutil.ReadFile(testcase.input)
		require.NoError(t, err)

		var actual CodePipelineCloudWatchEvent
		require.NoError(t, json.Unmarshal(data, &actual))

		require.Equal(t, testcase.expect, actual)
	}
}
