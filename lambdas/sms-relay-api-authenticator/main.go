package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhouziqunzzq/sms-relay-server/common"
)

const (
	defaultAWSRegion = "us-west-2"
	jwtSecretName    = "JWTSecret"
)

var (
	logger        = log.Default()
	secretsClient *secretsmanager.Client
)

type AuthRequest struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	IsAuthorized bool   `json:"isAuthorized"`
	Message      string `json:"message"`
}

func init() {
	awsRegion := defaultAWSRegion
	if region := os.Getenv("AWS_REGION"); region != "" {
		awsRegion = region
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		logger.Fatalf("unable to load SDK config, %v", err)
	}
	secretsClient = secretsmanager.NewFromConfig(cfg)
	logger.Println("Secrets Manager client initialized")
}

func handler(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// Extract the token from the Authorization header
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(request.AuthorizationToken, bearerPrefix) {
		logger.Println("invalid authorization header format")
		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID:    "unknown",
			PolicyDocument: generatePolicy("unknown", "Deny", request.MethodArn),
		}, nil
	}
	tokenString := strings.TrimPrefix(request.AuthorizationToken, bearerPrefix)

	// Retrieve the JWT secret using the helper function from the common package
	jwtSecretKey, err := common.GetSecretValue(ctx, secretsClient, jwtSecretName, "JWTKey")
	if err != nil {
		logger.Printf("failed to retrieve JWT secret: %v", err)
		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID:    "unknown",
			PolicyDocument: generatePolicy("unknown", "Deny", request.MethodArn),
		}, nil
	}

	// Parse and validate the JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		logger.Printf("invalid token: %v", err)
		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID:    "unknown",
			PolicyDocument: generatePolicy("unknown", "Deny", request.MethodArn),
		}, nil
	}

	// Extract claims and set PrincipalID
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Println("failed to parse token claims")
		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID:    "unknown",
			PolicyDocument: generatePolicy("unknown", "Deny", request.MethodArn),
		}, nil
	}

	logger.Printf("user %s authenticated successfully", claims["sub"])
	principalID, _ := claims["sub"].(string)
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:    principalID,
		PolicyDocument: generatePolicy(principalID, "Allow", request.MethodArn),
		Context: map[string]any{
			"user_id":   principalID,
			"user_type": claims["user_type"],
			"user_name": claims["user_name"],
			"device_id": claims["device_id"],
		},
	}, nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerPolicy {
	return events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   effect,
				Resource: []string{resource},
			},
		},
	}
}

func main() {
	lambda.Start(handler)
}
