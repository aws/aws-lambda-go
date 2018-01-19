# aws-lambda-go
Libraries, samples and tools to help Go developers develop AWS Lambda functions.

To learn more about writing AWS Lambda functions in Go, go to [the offical documentation](https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model.html)

# Getting Started

``` Go
// main.go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func hello() (string, error) {
	return "Hello ƛ!", nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
}
```
``` shell
# Remember to build your handler executable for linux!
GOOS=linux go build -o main main.go
```

# Deploying your functions

Take a look at the offical documentation for [deploying using the AWS CLI, AWS Cloudformation, and AWS SAM](https://docs.aws.amazon.com/lambda/latest/dg/deploying-lambda-apps.html)


# Event Integrations

If you're using AWS Lambda with an AWS event source, you can use one of the [event models](https://github.com/aws/aws-lambda-go/tree/master/events) for your request type.

[Check out the docs](https://docs.aws.amazon.com/lambda/latest/dg/use-cases.html) for detailed walkthroughs.

