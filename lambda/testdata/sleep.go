package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Event struct {
	SleepMilliseconds int `json:"sleep_ms"`
}

func handler(ctx context.Context, event Event) (string, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	logger := slog.Default().With("handler", "sleep-test")

	logger.Info("processing", "request_id", lc.AwsRequestID, "sleep_ms", event.SleepMilliseconds)
	time.Sleep(time.Duration(event.SleepMilliseconds) * time.Millisecond)
	logger.Info("completed", "request_id", lc.AwsRequestID)

	return "ok", nil
}

func main() {
	lambda.Start(handler)
}
