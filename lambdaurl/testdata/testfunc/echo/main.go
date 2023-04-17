package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(req any) (any, error) {
		return req, nil
	})
}
