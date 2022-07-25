package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type handler struct{}

func (h handler) Invoke(_ context.Context, e []byte) ([]byte, error) {
	return e, nil
}

func main() {
	lambda.Start(handler{})
}
