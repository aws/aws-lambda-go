package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalCodeBuildEvent(t *testing.T) {
	tests := []struct {
		input  string
		expect CodeBuildEvent
	}{
		{
			input: "testdata/codebuild-state-change.json",
			expect: CodeBuildEvent{
				Version:    "0",
				ID:         "c030038d-8c4d-6141-9545-00ff7b7153EX",
				DetailType: CodeBuildStateChangeDetailType,
				Source:     CodeBuildEventSource,
				AccountID:  "123456789012",
				Time:       time.Date(2017, 9, 1, 16, 14, 28, 0, time.UTC),
				Region:     "us-west-2",
				Resources: []string{
					"arn:aws:codebuild:us-west-2:123456789012:build/my-sample-project:8745a7a9-c340-456a-9166-edf953571bEX",
				},
				Detail: CodeBuildEventDetail{
					BuildStatus: CodeBuildPhaseStatusSucceeded,
					ProjectName: "my-sample-project",
					BuildID:     "arn:aws:codebuild:us-west-2:123456789012:build/my-sample-project:8745a7a9-c340-456a-9166-edf953571bEX",
					AdditionalInformation: CodeBuildEventAdditionalInformation{
						Artifact: CodeBuildArtifact{
							MD5Sum:    "da9c44c8a9a3cd4b443126e823168fEX",
							SHA256Sum: "6ccc2ae1df9d155ba83c597051611c42d60e09c6329dcb14a312cecc0a8e39EX",
							Location:  "arn:aws:s3:::codebuild-123456789012-output-bucket/my-output-artifact.zip",
						},
						Environment: CodeBuildEnvironment{
							Image:          "aws/codebuild/standard:2.0",
							PrivilegedMode: false,
							ComputeType:    "BUILD_GENERAL1_SMALL",
							Type:           "LINUX_CONTAINER",
							EnvironmentVariables: []CodeBuildEnvironmentVariable{
								{
									Name:  "TEST",
									Type:  "PLAINTEXT",
									Value: "TEST",
								},
							},
						},
						Timeout:        DurationMinutes(60 * time.Minute),
						BuildComplete:  true,
						Initiator:      "MyCodeBuildDemoUser",
						BuildStartTime: CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
						Source: CodeBuildSource{
							Location: "codebuild-123456789012-input-bucket/my-input-artifact.zip",
							Type:     "S3",
						},
						Logs: CodeBuildLogs{
							GroupName:  "/aws/codebuild/my-sample-project",
							StreamName: "8745a7a9-c340-456a-9166-edf953571bEX",
							DeepLink:   "https://console.aws.amazon.com/cloudwatch/home?region=us-west-2#logEvent:group=/aws/codebuild/my-sample-project;stream=8745a7a9-c340-456a-9166-edf953571bEX",
						},
						Phases: []CodeBuildPhase{
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeSubmitted,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2019, 9, 13, 4, 12, 29, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeQueued,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 05, 0, time.UTC)),
								Duration:     DurationSeconds(36 * time.Second),
								PhaseType:    CodeBuildPhaseTypeProvisioning,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 5, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								Duration:     DurationSeconds(4 * time.Second),
								PhaseType:    CodeBuildPhaseTypeDownloadSource,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeInstall,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypePreBuild,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								Duration:     DurationSeconds(70 * time.Second),
								PhaseType:    CodeBuildPhaseTypeBuild,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypePostBuild,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeUploadArtifacts,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 26, 0, time.UTC)),
								Duration:     DurationSeconds(4 * time.Second),
								PhaseType:    CodeBuildPhaseTypeFinalizing,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								StartTime: CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 26, 0, time.UTC)),
								PhaseType: CodeBuildPhaseTypeCompleted,
							},
						},
					},
					CurrentPhase:        CodeBuildPhaseTypeCompleted,
					CurrentPhaseContext: "[]",
					Version:             "1",
				},
			},
		},
		{
			input: "testdata/codebuild-phase-change.json",
			expect: CodeBuildEvent{
				Version:    "0",
				ID:         "43ddc2bd-af76-9ca5-2dc7-b695e15adeEX",
				DetailType: CodeBuildPhaseChangeDetailType,
				Source:     CodeBuildEventSource,
				AccountID:  "123456789012",
				Time:       time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC),
				Region:     "us-west-2",
				Resources: []string{
					"arn:aws:codebuild:us-west-2:123456789012:build/my-sample-project:8745a7a9-c340-456a-9166-edf953571bEX",
				},
				Detail: CodeBuildEventDetail{
					CompletedPhase:        CodeBuildPhaseTypeCompleted,
					ProjectName:           "my-sample-project",
					BuildID:               "arn:aws:codebuild:us-west-2:123456789012:build/my-sample-project:8745a7a9-c340-456a-9166-edf953571bEX",
					CompletedPhaseContext: "[]",
					AdditionalInformation: CodeBuildEventAdditionalInformation{
						Artifact: CodeBuildArtifact{
							MD5Sum:    "da9c44c8a9a3cd4b443126e823168fEX",
							SHA256Sum: "6ccc2ae1df9d155ba83c597051611c42d60e09c6329dcb14a312cecc0a8e39EX",
							Location:  "arn:aws:s3:::codebuild-123456789012-output-bucket/my-output-artifact.zip",
						},
						Environment: CodeBuildEnvironment{
							Image:                "aws/codebuild/standard:2.0",
							PrivilegedMode:       false,
							ComputeType:          "BUILD_GENERAL1_SMALL",
							Type:                 "LINUX_CONTAINER",
							EnvironmentVariables: []CodeBuildEnvironmentVariable{},
						},
						Timeout:        DurationMinutes(60 * time.Minute),
						BuildComplete:  true,
						Initiator:      "MyCodeBuildDemoUser",
						BuildStartTime: CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
						Source: CodeBuildSource{
							Location: "codebuild-123456789012-input-bucket/my-input-artifact.zip",
							Type:     "S3",
						},
						Logs: CodeBuildLogs{
							GroupName:  "/aws/codebuild/my-sample-project",
							StreamName: "8745a7a9-c340-456a-9166-edf953571bEX",
							DeepLink:   "https://console.aws.amazon.com/cloudwatch/home?region=us-west-2#logEvent:group=/aws/codebuild/my-sample-project;stream=8745a7a9-c340-456a-9166-edf953571bEX",
						},
						Phases: []CodeBuildPhase{
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeSubmitted,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2019, 9, 13, 4, 12, 29, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeQueued,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 12, 29, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 05, 0, time.UTC)),
								Duration:     DurationSeconds(36 * time.Second),
								PhaseType:    CodeBuildPhaseTypeProvisioning,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 5, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								Duration:     DurationSeconds(4 * time.Second),
								PhaseType:    CodeBuildPhaseTypeDownloadSource,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeInstall,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypePreBuild,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 13, 10, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								Duration:     DurationSeconds(70 * time.Second),
								PhaseType:    CodeBuildPhaseTypeBuild,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypePostBuild,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								Duration:     DurationSeconds(0),
								PhaseType:    CodeBuildPhaseTypeUploadArtifacts,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								PhaseContext: []interface{}{},
								StartTime:    CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
								EndTime:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 26, 0, time.UTC)),
								Duration:     DurationSeconds(4 * time.Second),
								PhaseType:    CodeBuildPhaseTypeFinalizing,
								PhaseStatus:  CodeBuildPhaseStatusSucceeded,
							},
							{
								StartTime: CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 26, 0, time.UTC)),
								PhaseType: CodeBuildPhaseTypeCompleted,
							},
						},
					},
					CompletedPhaseStatus:   CodeBuildPhaseStatusSucceeded,
					CompletedPhaseDuration: DurationSeconds(4 * time.Second),
					Version:                "1",
					CompletedPhaseStart:    CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 21, 0, time.UTC)),
					CompletedPhaseEnd:      CodeBuildTime(time.Date(2017, 9, 1, 16, 14, 26, 0, time.UTC)),
				},
			},
		},
	}

	for _, testcase := range tests {
		data, err := ioutil.ReadFile(testcase.input)
		require.NoError(t, err)

		var actual CodeBuildEvent
		require.NoError(t, json.Unmarshal(data, &actual))

		require.Equal(t, testcase.expect, actual)
	}
}
