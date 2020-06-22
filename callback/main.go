package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ngs/akerun-slack-notifier/akerun"
)

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req Request) (Response, error) {
	code := req.QueryStringParameters["code"]
	tokenResp, err := akerun.Config.Exchange(ctx, code)
	json, err := json.MarshalIndent(tokenResp, "", "  ")
	if err != nil {
		resp := Response{
			StatusCode: 400,
			Body:       err.Error(),
		}
		return resp, nil
	}
	err = akerun.UploadToken(tokenResp)
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
