package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"net/url"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	clientID := os.Getenv("AKERUN_CLIENT_ID")
	redirectURL := os.Getenv("AKERUN_REDIRECT_URL")

	if clientID == "" || redirectURL == "" {
		resp := Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Body:            `<html><head><title>Configuration Error</title></head><body>AKERUN_CLIENT_ID or AKERUN_REDIRECT_URL is not configured</body></html>`,
			Headers: map[string]string{
				"Content-Type": "text/html",
			},
		}
		return resp, nil
	}
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURL)
	params.Add("response_type", "code")
	authURL := "https://api.akerun.com/oauth/authorize/?" + params.Encode()

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
