package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	// conventional SIGTERM callback
	signaled := make(chan os.Signal, 1)
	signal.Notify(signaled, syscall.SIGTERM)
	go func() {
		<-signaled
		fmt.Println("I've been TERMINATED!")
	}()

}

func main() {
	// lambda option to enable sigterm, plus optional extra sigterm callbacks
	sigtermOption := lambda.WithEnableSIGTERM(func() {
		fmt.Println("Hello SIGTERM!")
	})
	handlerOptions := []lambda.Option{}
	if os.Getenv("ENABLE_SIGTERM") != "" {
		handlerOptions = append(handlerOptions, sigtermOption)
	}
	lambda.StartWithOptions(
		func(ctx context.Context) {
			deadline, _ := ctx.Deadline()
			<-time.After(time.Until(deadline) + time.Second)
			panic("unreachable line reached!")
		},
		handlerOptions...,
	)
}
