package events_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleECSContainerInstanceEvent() {
	lambda.Start(func(ctx context.Context, ecsEvent *events.ECSContainerInstanceEvent) {
		outputJSON, _ := json.MarshalIndent(ecsEvent, "", " ")
		fmt.Printf("Data = %s", outputJSON)
	})
}
