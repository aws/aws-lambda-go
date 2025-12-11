package events_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleChimeBotEvent() {
	lambda.Start(func(_ context.Context, chimeBotEvent *events.ChimeBotEvent) error {
		switch chimeBotEvent.EventType {
		case "Invite":
			fmt.Printf("Thanks for inviting me to this room %s\n", chimeBotEvent.Sender.SenderID)
			return nil
		case "Mention":
			fmt.Printf("Thanks for mentioning me %s\n", chimeBotEvent.Sender.SenderID)
			return nil
		case "Remove":
			fmt.Printf("I have been removed from %q by %q\n", chimeBotEvent.Discussion.DiscussionType, chimeBotEvent.Sender.SenderID)
			return nil
		default:
			return fmt.Errorf("event type %q is unsupported", chimeBotEvent.EventType)
		}
	})
}
