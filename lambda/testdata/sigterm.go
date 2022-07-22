package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.StartWithOptions(
		func(ctx context.Context) {
			deadline, _ := ctx.Deadline()
			fmt.Println("yeet my meat")
			<-time.After(time.Until(deadline) + time.Second)
			panic("unreachable line reached!")
		},
		lambda.WithEnableSIGTERM(func() {
			fmt.Println("Hello SIGTERM!")
		}),
	)
}
