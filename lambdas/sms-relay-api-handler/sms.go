package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/zhouziqunzzq/sms-relay-server/models"
)

type SMSRequest struct {
	PhoneNumber string `json:"phone_number"` // Phone number of the sender, in E.164 format
	Body        string `json:"body"`         // Content of the SMS message
}

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

	// Validate user type and device ID
	userType := request.RequestContext.Authorizer["user_type"].(string)
	deviceID := request.RequestContext.Authorizer["device_id"].(string)
	if userType != models.UserTypeDevice || deviceID == "" {
		logger.Printf("invalid user type or device ID: %s, %s", userType, deviceID)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid user type or device ID. Only devices can send SMS.",
		}, nil
	}

	// Validate and parse the request body
	var smsReq SMSRequest
	if err := json.Unmarshal([]byte(request.Body), &smsReq); err != nil {
		logger.Printf("failed to unmarshal request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}
	if smsReq.PhoneNumber == "" || smsReq.Body == "" {
		logger.Println("phone number or body is empty")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Phone number and body are required",
		}, nil
	}

	// Get Device by ID
	device, err := getDeviceByID(ctx, deviceID)
	if err != nil {
		logger.Printf("failed to get device by ID: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}
	if device == nil {
		logger.Println("device not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Device not found",
		}, nil
	}

	// Get Phone Number by number
	phoneNumber, err := getPhoneNumberByPhoneNumber(ctx, smsReq.PhoneNumber)
	if err != nil {
		logger.Printf("failed to get phone number by number: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}
	if phoneNumber == nil {
		logger.Println("phone number not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Phone number not found",
		}, nil
	}

	// Validate that the phone number is associated with the device
	phoneNumberIDSet := make(map[string]struct{})
	for _, num := range device.PhoneNumberIDs {
		phoneNumberIDSet[num] = struct{}{}
	}
	if _, ok := phoneNumberIDSet[phoneNumber.ID]; !ok {
		logger.Println("phone number is not associated with the device")
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Phone number is not associated with the device",
		}, nil
	}

	// Construct the SQS message
	smsRelayRequest := models.SMSRelayRequest{
		Device:      *device,
		PhoneNumber: *phoneNumber,
		SMS: models.SMS{
			From:          smsReq.PhoneNumber,
			Body:          smsReq.Body,
			PhoneNumberID: phoneNumber.ID,
		},
	}

	// Send the SMSRelayRequest to SQS
	messageBody, err := json.Marshal(smsRelayRequest)
	if err != nil {
		logger.Printf("failed to marshal SMSRelayRequest: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}
	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(sqsQueueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	if err != nil {
		logger.Printf("failed to send message to SQS: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to send message",
		}, nil
	}
	logger.Println("SMSRelayRequest successfully sent to SQS")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Message sent successfully",
	}, nil
}
