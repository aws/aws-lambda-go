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

# Building your function

Preparing a binary to deploy to AWS Lambda requires that it is compiled for Linux and placed into a .zip file.

## For developers on Linux and macOS
``` shell
# Remember to build your handler executable for Linux!
GOOS=linux go build -o main main.go
zip main.zip main
```

## For developers on Windows

Windows developers may have trouble producing a zip file that marks the binary as exectuable on Linux. To create a .zip that will work on AWS Lambda, the `build-lambda-zip` tool may be helpful.

Get the tool
``` shell
go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
```

Use the tool from your `GOPATH`. If you have a default installation of Go, the tool will be in `%USERPROFILE%\Go\bin`. 

in cmd.exe:
``` bat
set GOOS=linux
go build -o main main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main
```

in Powershell:
``` posh
$env:GOOS = "linux"
go build -o main main.go
~\Go\Bin\build-lambda-zip.exe -o main.zip main
```
# Deploying your functions

To deploy your function, refer to the offical documentation for [deploying using the AWS CLI, AWS Cloudformation, and AWS SAM](https://docs.aws.amazon.com/lambda/latest/dg/deploying-lambda-apps.html).

# Event Integrations

The [event models](https://github.com/aws/aws-lambda-go/tree/master/events) can be used to model AWS event sources. The offical documentation has [detailed walkthroughs](https://docs.aws.amazon.com/lambda/latest/dg/use-cases.html).

