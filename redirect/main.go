package main

import (
	"context"

	"github.com/ngs/akerun-slack-notifier/akerun"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	authURL := akerun.Config.AuthCodeURL("")
	resp := Response{
		StatusCode:      302,
		IsBase64Encoded: false,
		Body:            `<html><head><title>Redirecting</title></head><body>Redirecting to <a href="` + authURL + `">` + authURL + `</a></body></html>`,
		Headers: map[string]string{
			"Content-Type": "text/html",
			"Location":     authURL,
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
