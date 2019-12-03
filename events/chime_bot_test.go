// Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestChimeBotEventMarshaling(t *testing.T) {
	// From https://docs.aws.amazon.com/chime/latest/ag/manage-chat-bots.html
	tests := map[string]struct {
		inputJSON     string
		expectedEvent ChimeBotEvent
	}{
		"Example Invite Event": {
			inputJSON: ` {
  "Sender": {
    "SenderId": "user@example.com",
    "SenderIdType": "EmailId"
  },
  "Discussion": {
    "DiscussionId": "abcdef12-g34h-56i7-j8kl-mn9opqr012st",
    "DiscussionType": "Room"
  },
  "EventType": "Invite",
  "InboundHttpsEndpoint": {
    "EndpointType": "Persistent",
    "Url": "https://hooks.a.chime.aws/incomingwebhooks/a1b2c34d-5678-90e1-f23g-h45i67j8901k?token=ABCDefGHiJK1LMnoP2Q3RST4uvwxYZAbC56DeFghIJkLM7N8OP9QRsTuV0WXYZABcdefgHiJ"
  },
  "EventTimestamp": "2019-04-04T21:27:52.736Z"
}`,
			expectedEvent: ChimeBotEvent{
				Sender: ChimeBotEventSender{
					SenderID:     "user@example.com",
					SenderIDType: "EmailId",
				},
				Discussion: ChimeBotEventDiscussion{
					DiscussionID:   "abcdef12-g34h-56i7-j8kl-mn9opqr012st",
					DiscussionType: "Room",
				},
				EventType: "Invite",
				InboundHTTPSEndpoint: &ChimeBotEventInboundHTTPSEndpoint{
					EndpointType: "Persistent",
					URL:          "https://hooks.a.chime.aws/incomingwebhooks/a1b2c34d-5678-90e1-f23g-h45i67j8901k?token=ABCDefGHiJK1LMnoP2Q3RST4uvwxYZAbC56DeFghIJkLM7N8OP9QRsTuV0WXYZABcdefgHiJ",
				},
				EventTimestamp: time.Date(2019, 04, 04, 21, 27, 52, 736000000, time.UTC),
				Message:        "",
			},
		},
		"Example Mention Event": {
			inputJSON: `{
  "Sender": {
    "SenderId": "user@example.com",
    "SenderIdType": "EmailId"
  },
  "Discussion": {
    "DiscussionId": "abcdef12-g34h-56i7-j8kl-mn9opqr012st",
    "DiscussionType": "Room"
  },
  "EventType": "Mention",
  "InboundHttpsEndpoint": {
    "EndpointType": "ShortLived",
    "Url": "https://hooks.a.chime.aws/incomingwebhooks/a1b2c34d-5678-90e1-f23g-h45i67j8901k?token=ABCDefGHiJK1LMnoP2Q3RST4uvwxYZAbC56DeFghIJkLM7N8OP9QRsTuV0WXYZABcdefgHiJ"
  },
  "EventTimestamp": "2019-04-04T21:30:43.181Z",
  "Message": "@botDisplayName@example.com Hello Chatbot"
}`,
			expectedEvent: ChimeBotEvent{
				Sender: ChimeBotEventSender{
					SenderID:     "user@example.com",
					SenderIDType: "EmailId",
				},
				Discussion: ChimeBotEventDiscussion{
					DiscussionID:   "abcdef12-g34h-56i7-j8kl-mn9opqr012st",
					DiscussionType: "Room",
				},
				EventType: "Mention",
				InboundHTTPSEndpoint: &ChimeBotEventInboundHTTPSEndpoint{
					EndpointType: "ShortLived",
					URL:          "https://hooks.a.chime.aws/incomingwebhooks/a1b2c34d-5678-90e1-f23g-h45i67j8901k?token=ABCDefGHiJK1LMnoP2Q3RST4uvwxYZAbC56DeFghIJkLM7N8OP9QRsTuV0WXYZABcdefgHiJ",
				},
				EventTimestamp: time.Date(2019, 04, 04, 21, 30, 43, 181000000, time.UTC),
				Message:        "@botDisplayName@example.com Hello Chatbot",
			},
		},
		"Example Remove Event": {
			inputJSON: `{
  "Sender": {
    "SenderId": "user@example.com",
    "SenderIdType": "EmailId"
  },
  "Discussion": {
    "DiscussionId": "abcdef12-g34h-56i7-j8kl-mn9opqr012st",
    "DiscussionType": "Room"
  },
  "EventType": "Remove",
  "EventTimestamp": "2019-04-04T21:27:29.626Z"
}`,
			expectedEvent: ChimeBotEvent{
				Sender: ChimeBotEventSender{
					SenderID:     "user@example.com",
					SenderIDType: "EmailId",
				},
				Discussion: ChimeBotEventDiscussion{
					DiscussionID:   "abcdef12-g34h-56i7-j8kl-mn9opqr012st",
					DiscussionType: "Room",
				},
				EventType:      "Remove",
				EventTimestamp: time.Date(2019, 04, 04, 21, 27, 29, 626000000, time.UTC),
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			var testEvent ChimeBotEvent
			if err := json.Unmarshal([]byte(test.inputJSON), &testEvent); err != nil {
				t.Errorf("could not unmarshal event. details: %v", err)
			}

			assert.Equal(t, testEvent, test.expectedEvent)

			outputJSON, err := json.Marshal(testEvent)
			if err != nil {
				t.Errorf("could not marshal event. details: %v", err)
			}
			assert.JSONEq(t, test.inputJSON, string(outputJSON))
		})
	}
}

func TestChimeBotMarshalingMalformedJSON(t *testing.T) {
	test.TestMalformedJson(t, ChimeBotEvent{})
}
