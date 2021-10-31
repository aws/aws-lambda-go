# Sample Function

The following is a sample class and Lambda function that receives a Custome Resource event. Please see the official doc [here](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources.html).

```go
import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func handler(ctx context.Context, customResourceEvent events.CustomResource) {
	fmt.Printf("The RequestType is %s\n", customResourceEvent.RequestType)

	switch inputEvent.RequestType {
	case "Update":
        // some work
	case "Delete":
        // some work
	case "Create":
        // some work
	}
}
```
