# AWS Lambda for Go 

[![tests](https://github.com/aws/aws-lambda-go/workflows/tests/badge.svg)](https://github.com/aws/aws-lambda-go/actions?query=workflow%3Atests)
[![Go Reference](https://pkg.go.dev/badge/github.com/aws/aws-lambda-go.svg)](https://pkg.go.dev/github.com/aws/aws-lambda-go)
[![GoCard](https://goreportcard.com/badge/github.com/aws/aws-lambda-go)](https://goreportcard.com/report/github.com/aws/aws-lambda-go)
[![codecov](https://codecov.io/gh/aws/aws-lambda-go/branch/master/graph/badge.svg)](https://codecov.io/gh/aws/aws-lambda-go)

Libraries, samples, and tools to help Go developers develop AWS Lambda functions.

To learn more about writing AWS Lambda functions in Go, go to [the official documentation](https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model.html)

# Getting Started

``` Go
// main.go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func hello() (string, error) {
	return "Hello λ!", nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
}
```

# Building your function

Preparing a binary to deploy to AWS Lambda requires that it is compiled for Linux and placed into a .zip file. When using the `provided`, `provided.al2`, or `provided.al2023` runtime, the executable within the .zip file should be named `bootstrap`. Lambda's default architecture is `x86_64`, so when cross compiling from a non-x86 environment, the executable should be built with `GOARCH=amd64`. Likewise, if the Lambda function will be [configured to use ARM](https://docs.aws.amazon.com/lambda/latest/dg/foundation-arch.html), the executable should built with `GOARCH=arm64`.

``` shell
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip lambda-handler.zip bootstrap
```

## For developers on Linux

On Linux, the Go compiler's default behavior is to link the output executable to the system libc for some standard library functionality (for example, DNS lookups). If the build environment is using a Linux distribution with a GNU libc version newer than the deployment environment, the application when deployed to Lambda may fail with an error like ``/lib64/libc.so.6: version `GLIBC_X.YZ' not found``. 

Most Go applications do not require linking to the system libc. This behavior can be disabled by using the `CGO_ENABLED` environment variable.

```
CGO_ENABLED=0 go build -o bootstrap main.go
zip lambda-handler.zip bootstrap
```

See [Using CGO](#using-cgo)

## For developers on Windows

Windows developers may have trouble producing a zip file that marks the binary as executable on Linux. To create a .zip that will work on AWS Lambda, the `build-lambda-zip` tool may be helpful.

Get the tool
``` shell
go.exe install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest
```

Use the tool from your `GOPATH`. If you have a default installation of Go, the tool will be in `%USERPROFILE%\Go\bin`. 

in cmd.exe:
``` bat
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o bootstrap main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o lambda-handler.zip bootstrap
```

in Powershell:
``` posh
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o bootstrap main.go
~\Go\Bin\build-lambda-zip.exe -o lambda-handler.zip bootstrap
```

## Using CGO

For applications that require CGO, the build environment must be using a GNU libc version installed compatible with the target Lambda runtime. Otherwise, execution may fail with errors like ``/lib64/libc.so.6: version `GLIBC_X.YZ' not found``.

| Lambda runtime  | GLIBC version
| ----- | ---
| `provided.al2023` | 2.34
| `provided.al2` | 2.26
| `provided` and `go1.x` | 2.17


Alternatively, Lambda supports container images as a deployment package alternative to .zip files. For more information, refer to the official documentation for [working with with container images](https://docs.aws.amazon.com/lambda/latest/dg/images-create.html).

# Deploying your functions

To deploy your function, refer to the official documentation for [deploying using the AWS CLI, AWS Cloudformation, and AWS SAM](https://docs.aws.amazon.com/lambda/latest/dg/deploying-lambda-apps.html).

# Event Integrations

[![events package](https://pkg.go.dev/badge/github.com/aws/aws-lambda-go/events.svg)](https://pkg.go.dev/github.com/aws/aws-lambda-go/events)

The events package provides type definitions for common AWS event sources. The official documentation has [detailed walkthroughs](https://docs.aws.amazon.com/lambda/latest/dg/use-cases.html).

