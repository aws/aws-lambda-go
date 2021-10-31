# Sample Function

The following is a sample class and Lambda function that receives a CloudFormation Custom Resource event. Please see the official doc [here](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources.html).

```go
import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func handler(ctx context.Context, cfnCustomResourceEvent events.CfnCustomResource) {
	fmt.Printf("The RequestType is %s\n", cfnCustomResourceEvent.RequestType)

	switch cfnCustomResourceEvent.RequestType {
	case "Update":
        // some work
	case "Delete":
        // some work
	case "Create":
        // some work
	}
}
```
