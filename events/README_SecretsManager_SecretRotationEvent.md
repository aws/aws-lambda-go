# Sample Function

The following is a sample Lambda function that handles a SecretsManager secret rotation event.

```go
package main

import (
    "fmt"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
)

func handler(ctx context.Context, event SecretsManagerSecretRotationEvent) error {
	fmt.Printf("rotating secret %s with token %s\n", 
        event.SecretID, event.ClientRequestToken)

    switch event.Step {
	case "createSecret":
		// create
	case "setSecret":
		// set
	case "finishSecret":
		// finish
	case "testSecret":
		// test
	}

    return nil
}


func main() {
	lambda.Start(handler)
}
```