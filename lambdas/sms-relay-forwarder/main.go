package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/zhouziqunzzq/sms-relay-server/models"
)

const (
	smtpUsernameSecretName = "SMTPUsername"
	smtpPasswordSecretName = "SMTPPassword"
)

var (
	logger = log.Default()

	smtpServer    string
	smtpPort      string
	secretsClient *secretsmanager.Client
	useSSL        bool
)

func init() {
	// Load SMTP server address and port from environment variables
	smtpServer = os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		log.Fatalf("SMTP_SERVER environment variable is not set")
	}

	smtpPort = os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		log.Fatalf("SMTP_PORT environment variable is not set")
	}

	// Load SSL flag from environment variable
	sslFlag := os.Getenv("SSL")
	if sslFlag == "true" {
		useSSL = true
	} else {
		useSSL = false
	}

	// Initialize AWS Secrets Manager client
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	secretsClient = secretsmanager.NewFromConfig(cfg)
	log.Println("Secrets Manager client initialized")
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var smsRelayRequest models.SMSRelayRequest
		if err := json.Unmarshal([]byte(message.Body), &smsRelayRequest); err != nil {
			logger.Printf("failed to unmarshal SQS message: %v", err)
			return err
		}

		logger.Printf("Processing SMS Relay Request (Device ID: %s, Phone Number: %s)",
			smsRelayRequest.Device.ID, smsRelayRequest.PhoneNumber.PhoneNumber)

		// Forward SMS by Email
		if err := forwardSMSByEmail(ctx, smsRelayRequest); err != nil {
			logger.Printf("failed to forward SMS by email: %v", err)
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
