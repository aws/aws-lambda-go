# Overview

This package provides input types for Lambda functions that process Amazon Cognito events.

# Sample Function

The following is a sample Lambda function that receives Amazon Cognito event record data as an input and writes some of the record data to CloudWatch Logs. (Note that by default anything written to Console will be logged as CloudWatch Logs events.)

```go

import (
    "strings"
    "github.com/aws/aws-lambda-go/events/cognitoevents"
)

func handleRequest(ctx context.Context, cognitoEvent cognitoevents.CognitoEvent) {
    for datasetName, datasetRecord := range e.DatasetRecords {
        fmt.Printf("[%s -- %s] %s -> %s -> %s \n",
            cognitoEvent.EventType,
            datasetName,
            datasetRecord.OldValue,
            datasetRecord.Op,
            datasetRecord.NewValue)
    }
}
```
