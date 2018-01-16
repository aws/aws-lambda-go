# aws-lambda-go
Libraries, samples and tools to help Go developers develop AWS Lambda functions.

- [AWS Lambda Go](#aws-lambda-go)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Simple Handler That Responds to POST message](#simple-handler-that-responds-to-post-message)
    - [Sample Input Types for Lambda functions that process AWS events](#sample-input-types-for-lambda-functions-that-process-aws-events)
  - [License](#license)

## Installation

You can install the library to your GOPATH using `go get`.

```
go get github.com/aws/aws-lambda-go ./...
```

## Usage

Please see the following blog post by Chris Munns for full tutorial on how to use the library: [Announcing Go Support for AWS Lambda](https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/)

You can also see official AWS examples of using AWS Lambda with Go [here](https://github.com/aws-samples/lambda-go-samples)

### Simple Handler That Responds to POST message

```go
package main

import (
 "errors"
 "log"
 "github.com/aws/aws-lambda-go/events"
 "github.com/aws/aws-lambda-go/lambda"
)

var (
 // ErrNameNotProvided is thrown when a name is not provided
 ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

 // stdout and stderr are sent to AWS CloudWatch Logs
 log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

 // If no name is provided in the HTTP request body, throw an error
 if len(request.Body) < 1 {
  return events.APIGatewayProxyResponse{}, ErrNameNotProvided
 }

 return events.APIGatewayProxyResponse{
  Body:       "Hello " + request.Body,
  StatusCode: 200,
 }, nil

}

func main() {
 lambda.Start(Handler)
}
```

This can be packaged and deployed as an AWS Lambda function with the following commands. (This assumes that you have already configured [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/installing.html) on your machine.

```bash
# build go binary and zip 
$ GOOS=linux go build -o main
$ zip deployment.zip main

# create new AWS Lambda function and API handler using SAM
$ cat > template.yml <<EOF
AWSTemplateFormatVersion: 2010-09-09
Transform: AWS::Serverless-2016-10-31
Resources:
  HelloFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: main
      Runtime: go1.x
      Tracing: Active
      Events:
        GetEvent:
          Type: Api
          Properties:
            Path: /
            Method: post
EOF

# package function with cloudformation
$ aws cloudformation package --template-file template.yml --s3-bucket <<S3_BUCKET>> --s3-prefix <<S3_DIR>> --output-template-file packaged-template.yml

# deploy function with cloudformation
$ aws cloudformation deploy --template-file packaged-template.yml --stack-name helloworld --capabilities CAPABILITY_IAM

# Test function with curl you can get the api-gateway-url from the AWS Console
$ curl -XPOST -d "Werner" https://<<api-gateway-url>>
Hello Werner
```

Check out the [AWS announcement blog post](https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/) for more information and further reading.

## Sample Input Types for Lambda functions that process AWS events
- [Overview](events/README.md)
  - [API Gateway Events](events/README_ApiGatewayEvent.md)
  - [Cognito Events](events/README_Cognito.md)
  - [Config Events](events/README_Config.md)
  - [DynamoDB Events](events/README_DynamoDB.md)
  - [Kinesis Events](events/README_Kinesis.md)
  - [Kinesis Firehose Events](events/README_KinesisFirehose.md)
  - [S3 Events](events/README_S3.md)
  - [SNS Events](events/README_SNS.md)

## License

Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Lambda functions are made available under a modified MIT license.
See LICENSE-LAMBDACODE for details.

The remainder of the project is made available under the terms of the
Apache License, version 2.0. See LICENSE for details.

