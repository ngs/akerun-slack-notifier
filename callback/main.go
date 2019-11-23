package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ngs/akerun-slack-notifier/akerun"
)

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req Request) (Response, error) {
	code := req.QueryStringParameters["code"]
	bucket := os.Getenv("S3_BUCKET")
	tokenResp, err := akerun.Config.Exchange(ctx, code)
	if err != nil {
		resp := Response{
			StatusCode: 400,
			Body:       err.Error(),
		}
		return resp, nil
	}
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)
	json, err := json.MarshalIndent(tokenResp, "", "  ")
	if err != nil {
		resp := Response{
			StatusCode: 400,
			Body:       err.Error(),
		}
		return resp, nil
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("token.json"),
		Body:   bytes.NewReader(json),
	})
	if err != nil {
		resp := Response{
			StatusCode: 400,
			Body:       err.Error(),
		}
		return resp, nil
	}
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(json),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
