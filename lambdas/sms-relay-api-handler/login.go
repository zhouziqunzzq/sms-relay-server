package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/zhouziqunzzq/sms-relay-server/common"
	"github.com/zhouziqunzzq/sms-relay-server/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User             *models.User `json:"user,omitempty"`
	Token            string       `json:"token,omitempty"`
	TokenExpireAfter string       `json:"token_expire_after,omitempty"`
}

func handlePostLogin(ctx context.Context, request events.APIGatewayProxyRequest) (
	resp events.APIGatewayProxyResponse, err error,
) {
	// Being defensive here - API Gateway should've already filtered out non-POST requests
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "Method Not Allowed",
		}, nil
	}

	// Validate and parse the request body
	var loginReq LoginRequest
	if err := json.Unmarshal([]byte(request.Body), &loginReq); err != nil {
		logger.Printf("error unmarshalling login request: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}
	if loginReq.Username == "" || loginReq.Password == "" {
		logger.Println("username or password is empty")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Username and password are required",
		}, nil
	}

	// Fetch user from DynamoDB
	user, err := getUserByUsername(ctx, loginReq.Username)
	if err != nil {
		logger.Println("error fetching user")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}
	if user == nil {
		logger.Println("user not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Username or password is incorrect",
		}, nil
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			logger.Println("password not matched")
		} else {
			logger.Println("error validating password")
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Username or password is incorrect",
		}, nil
	}

	// Fetch the secret value for JWT
	jwtSigningKey, err := common.GetSecretValue(ctx, secretsClient, jwtSecretName, "JWTKey")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	// Generate JWT token
	expirationTime := time.Now().Add(jwtValidityDuration)
	signedToken, err := user.GenerateJWT([]byte(jwtSigningKey), expirationTime)
	if err != nil {
		logger.Println("error generating JWT token")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	// Generate and return the response
	response := LoginResponse{
		User:             user,
		Token:            signedToken,
		TokenExpireAfter: expirationTime.Format(time.RFC3339),
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Printf("error marshalling response: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}
	logger.Printf("user %s logged in successfully\n", user.Username)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}
