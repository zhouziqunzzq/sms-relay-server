package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func handleUser(ctx context.Context, request events.APIGatewayProxyRequest) (
	resp events.APIGatewayProxyResponse, err error,
) {
	switch request.HTTPMethod {
	case "GET":
		return handleGetUser(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "Method Not Allowed",
		}, nil
	}
}

// handleGetUser retrieves user information based on the user ID provided in the authorization context.
// It returns the user details in the response body.
func handleGetUser(ctx context.Context, request events.APIGatewayProxyRequest) (
	resp events.APIGatewayProxyResponse, err error,
) {
	// Extract user ID from the request context
	userID, ok := request.RequestContext.Authorizer["user_id"].(string)
	if !ok || userID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "User ID not found in authorization context",
		}, nil
	}

	// Fetch user details from the database
	user, err := getUserByID(ctx, userID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to retrieve user details",
		}, err
	}
	if user == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "User not found",
		}, nil
	}

	// Return user details in the response
	responseBody, err := json.Marshal(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal user details",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}
