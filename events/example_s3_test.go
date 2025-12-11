package events_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ExampleS3Event() {
	lambda.Start(func(ctx context.Context, s3Event *events.S3Event) {
		for _, record := range s3Event.Records {
			s3 := record.S3
			fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
		}
	})
}

// S3 Object Lambda allows you to add your own code to S3 GET requests to modify and process data
// as it is returned to an application. This example receives S3 Object Lambda event data and returns object metadata.
func ExampleS3ObjectLambdaEvent() {
	lambda.Start(func(ctx context.Context, event *events.S3ObjectLambdaEvent) error {
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

		type Metadata struct {
			Length int
			Md5    string
		}
		type TransformedObject struct {
			Metadata Metadata
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

		// To complete the example, use the AWS SDK to write the response:
		//
		// import (
		// 	"github.com/aws/aws-sdk-go-v2/config"
		// 	"github.com/aws/aws-sdk-go-v2/service/s3"
		// )
		//
		// cfg, err := config.LoadDefaultConfig(context.TODO())
		// if err != nil {
		// 	return err
		// }
		// svc := s3.NewFromConfig(cfg)
		// input := &s3.WriteGetObjectResponseInput{
		// 	RequestRoute: &event.GetObjectContext.OutputRoute,
		// 	RequestToken: &event.GetObjectContext.OutputToken,
		// 	Body:         strings.NewReader(string(jsonData)),
		// }
		// res, err := svc.WriteGetObjectResponse(ctx, input)
		// if err != nil {
		// 	return err
		// }
		// fmt.Printf("%v", res)

		fmt.Printf("Transformed object metadata: %s\n", jsonData)
		return nil
	})
}

func toMd5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
