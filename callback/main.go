package main

import (
	"context"

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
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "OK",
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
