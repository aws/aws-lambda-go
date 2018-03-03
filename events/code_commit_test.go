package events

import (
	"testing"

	"github.com/aws/aws-lambda-go/events/test"
)

func TestCodeCommitReference(t *testing.T) {
	cases := []struct {
		Name  string
		Input []byte
	}{
		{
			Name: "CodeCommitReference",
			Input: []byte(`
        {
          "commit": "5c4ef1049f1d27deadbeeff313e0730018be182b",
          "ref": "refs/heads/master"
        }
      `),
		},
		{
			Name: "Created CodeCommitReference",
			Input: []byte(`
        {
          "commit": "5c4ef1049f1d27deadbeeff313e0730018be182b",
          "ref": "refs/heads/master",
          "created": true
        }
      `),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			test.AssertJsonBytes(t, c.Input, &CodeCommitReference{})
		})
	}
}

func TestCodeCommitCodeCommit(t *testing.T) {
	cases := []struct {
		Name  string
		Input []byte
	}{
		{
			Name:  "Empty CodeCommitReferences",
			Input: []byte(`{"references": []}`),
		},
		{
			Name: "CodeCommitCodeCommit",
			Input: []byte(`
        {
          "references": [
            {
              "commit": "5c4ef1049f1d27deadbeeff313e0730018be182b",
              "ref": "refs/heads/master"
            },
            {
              "commit": "5c4ef1049f1d27deadbeeff313e0730018be182b",
              "ref": "refs/heads/master",
              "created": true
            }
          ]
        }
      `),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			test.AssertJsonBytes(t, c.Input, &CodeCommitCodeCommit{})
		})
	}
}

func TestCodeCommitRecord(t *testing.T) {
	cases := []struct {
		Name  string
		Input []byte
	}{
		{
			Name: "CodeCommitRecord",
			Input: []byte(`
        {
          "eventId": "5a824061-17ca-46a9-bbf9-114edeadbeef",
          "eventVersion": "1.0",
          "eventTime": "2018-01-22T15:58:33.475+0000",
          "eventTriggerName": "my-trigger",
          "eventPartNumber": 1,
          "codecommit": {
            "references": []
          },
          "eventName": "TriggerEventTest",
          "eventTriggerConfigId": "5a824061-17ca-46a9-bbf9-114edeadbeef",
          "eventSourceARN": "arn:aws:codecommit:us-east-1:123456789012:my-repo",
          "userIdentityARN": "arn:aws:iam::123456789012:root",
          "eventSource": "aws:codecommit",
          "awsRegion": "us-east-1",
          "eventTotalParts": 1
        }
      `),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			test.AssertJsonBytes(t, c.Input, &CodeCommitRecord{})
		})
	}
}

func TestCodeCommitEventFile(t *testing.T) {
	test.AssertJsonFile(t, "./testdata/code-commit-event.json", &CodeCommitEvent{})
}

func TestCodeCommitEvent(t *testing.T) {
	cases := []struct {
		Name  string
		Input []byte
	}{
		{
			Name:  "Empty CodeCommitRecord",
			Input: []byte(`{"Records": []}`),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			test.AssertJsonBytes(t, c.Input, &CodeCommitEvent{})
		})
	}
}
