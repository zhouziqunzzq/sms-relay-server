package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	defaultAWSRegion = "us-west-2"

	userTableName        = "UserTable"
	usernameIndexName    = "UsernameIndex"
	deviceTableName      = "DeviceTable"
	phoneNumberTableName = "PhoneNumberTable"
	phoneNumberIndexName = "PhoneNumberIndex"

	jwtSecretName       = "JWTSecret"
	jwtValidityDuration = time.Hour * 24 * 7 // 7 days
)

var (
	logger        = log.Default()
	dbClient      *dynamodb.Client
	secretsClient *secretsmanager.Client
	sqsClient     *sqs.Client
	sqsQueueURL   string
)

// init initializes the DynamoDB and Secrets Manager clients.
func init() {
	// Initialize AWS clients
	awsRegion := defaultAWSRegion
	if region := os.Getenv("AWS_REGION"); region != "" {
		awsRegion = region
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		logger.Fatalf("unable to load SDK config, %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
	secretsClient = secretsmanager.NewFromConfig(cfg)
	sqsClient = sqs.NewFromConfig(cfg)
	logger.Println("DynamoDB, Secrets Manager, and SQS clients initialized")

	// Get the SQS queue URL from the environment variable
	sqsQueueURL = os.Getenv("SMS_RELAY_REQUEST_QUEUE_URL")
	if sqsQueueURL == "" {
		logger.Fatalf("SMS_RELAY_REQUEST_QUEUE_URL environment variable is not set")
	}
}

// handler processes incoming API Gateway requests and routes them to the appropriate function
// based on the request path.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.Path {
	case "/login":
		return handlePostLogin(ctx, request)
	case "/sms":
		return handlePostSMS(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
