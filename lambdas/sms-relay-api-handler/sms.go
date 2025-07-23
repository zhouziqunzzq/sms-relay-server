package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func handlePostSMS(ctx context.Context, request events.APIGatewayProxyRequest) (
	resp events.APIGatewayProxyResponse, err error,
) {
	// Being defensive here - API Gateway should've already filtered out non-POST requests
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "Method Not Allowed",
		}, nil
	}

	// Log requests for now
	logger.Printf("received request: %+v\n", request)
	logger.Printf("request authorizer context: %+v\n", request.RequestContext.Authorizer)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "SMS Relay API Handler",
	}, nil
}
