# Overview

CloudFormation has a different way of responding to most events due to the way stacks execute.

It is best to catch all errors and ensure the correct response is sent to the pre-signed URL that comes with the event.

To make this easier, a wrapper exists to allow the creation of custom resources without having to handle that.

# Sample Function

This sample will safely 'Echo' back anything given into the Echo parameter within the Custom Resource call.

```go

import (
    "context"
    "fmt"

    "github.com/aws/aws-lambda-go/cfn"
    "github.com/aws/aws-lambda-go/lambda"
)

func echoResource(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
    v, _ := event.ResourceProperties["Echo"].(string)

    data = map[string]interface{} {
        "Echo": v,
    }

    return
}

func main() {
	lambda.Start(cfn.LambdaWrap(echoResource))
}
```
