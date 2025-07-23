package common

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// GetSecretValue retrieves the value of a secret from AWS Secrets Manager.
func GetSecretValue(ctx context.Context, secretsClient *secretsmanager.Client, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}
	result, err := secretsClient.GetSecretValue(ctx, input)
	if err != nil {
		var resourceNotFound *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFound) {
			log.Printf("secret %s not found\n", secretName)
			return "", err
		}
		log.Printf("error fetching secret %s: %v\n", secretName, err)
		return "", err
	}
	return *result.SecretString, nil
}
