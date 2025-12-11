// Package lambda provides the runtime interface for AWS Lambda functions written in Go.
//
// This package handles communication with the Lambda execution environment, including
// receiving invocation events, managing context, and returning responses.
//
// Use lambda.Start() to register your handler function and begin processing Lambda invocations.
//
// See https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
package lambda
