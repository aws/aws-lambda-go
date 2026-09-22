package events_test

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func ExampleLexV2Event() {
	event := events.LexV2Event{
		MessageVersion:   "1.0",
		InvocationSource: "FulfillmentCodeHook",
		InputMode:        "Text",
		SessionID:        "12345678-1234-1234-1234-123456789012",
		InputTranscript:  "check my balance",
		Bot: events.LexV2Bot{
			ID:       "BOTID",
			Name:     "BankingBot",
			AliasID:  "ALIASID",
			LocaleID: "en_US",
			Version:  "1",
		},
		SessionState: events.LexV2SessionState{
			Intent: &events.LexV2Intent{
				Name:              "CheckBalance",
				State:             "ReadyForFulfillment",
				ConfirmationState: "None",
				Slots:             map[string]events.LexV2Slot{},
			},
		},
	}

	fmt.Printf("Bot: %s, Intent: %s\n", event.Bot.Name, event.SessionState.Intent.Name)
	// Output: Bot: BankingBot, Intent: CheckBalance
}

func ExampleLexV2Response() {
	response := events.LexV2Response{
		SessionState: events.LexV2SessionState{
			DialogAction: &events.LexV2DialogAction{
				Type: "Close",
			},
			Intent: &events.LexV2Intent{
				Name:              "CheckBalance",
				State:             "Fulfilled",
				ConfirmationState: "None",
				Slots:             map[string]events.LexV2Slot{},
			},
		},
		Messages: []events.LexV2Message{
			{
				ContentType: "PlainText",
				Content:     "Your balance is $1,234.56",
			},
		},
	}

	fmt.Printf("Action: %s, Message: %s\n", response.SessionState.DialogAction.Type, response.Messages[0].Content)
	// Output: Action: Close, Message: Your balance is $1,234.56
}
