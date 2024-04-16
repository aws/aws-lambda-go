# Sample Function

The following is a sample class and Lambda function that receives Amazon S3 Object Lambda event record data as an input and returns an object metadata output.

```go

// main.go
package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func handler(ctx context.Context, event events.S3ObjectLambdaEvent) error {
	url := event.GetObjectContext.InputS3URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	transformedObject := TransformedObject{
		Metadata: Metadata{
			Length: len(bodyBytes),
			Md5:    toMd5(bodyBytes),
		},
	}
	jsonData, err := json.Marshal(transformedObject)
	if err != nil {
		return err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	svc := s3.NewFromConfig(cfg)
	input := &s3.WriteGetObjectResponseInput{
		RequestRoute: &event.GetObjectContext.OutputRoute,
		RequestToken: &event.GetObjectContext.OutputToken,
		Body:         strings.NewReader(string(jsonData)),
	}
	res, err := svc.WriteGetObjectResponse(ctx, input)
	if err != nil {
		return err
	}
	fmt.Printf("%v", res)
	return nil
}

func toMd5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}

type TransformedObject struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Length int    `json:"length"`
	Md5    string `json:"md5"`
}

func main() {
	lambda.Start(handler)
}

```
