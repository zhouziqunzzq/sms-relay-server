package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/zhouziqunzzq/sms-relay-server/models"
)

func getUserByUsername(ctx context.Context, username string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(userTableName),
		IndexName:              aws.String(usernameIndexName),
		KeyConditionExpression: aws.String("Username = :username"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: username},
		},
	}

	result, err := dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil // User not found
	}

	var user models.User
	if err := attributevalue.UnmarshalMap(result.Items[0], &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func getDeviceByID(ctx context.Context, deviceID string) (*models.Device, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(deviceTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: deviceID},
		},
	}

	result, err := dbClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil // Device not found
	}

	var device models.Device
	if err := attributevalue.UnmarshalMap(result.Item, &device); err != nil {
		return nil, err
	}

	return &device, nil
}

func getPhoneNumberByPhoneNumber(ctx context.Context, number string) (*models.PhoneNumber, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(phoneNumberTableName),
		IndexName:              aws.String(phoneNumberIndexName),
		KeyConditionExpression: aws.String("PhoneNumber = :phoneNumber"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":phoneNumber": &types.AttributeValueMemberS{Value: number},
		},
	}

	result, err := dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil // Phone number not found
	}

	var phoneNumber models.PhoneNumber
	if err := attributevalue.UnmarshalMap(result.Items[0], &phoneNumber); err != nil {
		return nil, err
	}

	return &phoneNumber, nil
}
