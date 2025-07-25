package common

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// GetSecretValue retrieves the value of a secret from AWS Secrets Manager. If a key is provided,
// it parses the JSON secret and returns the value for the key.
func GetSecretValue(ctx context.Context, secretsClient *secretsmanager.Client, secretName string, key string) (string, error) {
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

	if key == "" {
		return *result.SecretString, nil
	}

	var secretMap map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretMap); err != nil {
		log.Printf("error parsing secret %s as JSON: %v\n", secretName, err)
		return "", err
	}

	value, exists := secretMap[key]
	if !exists {
		log.Printf("key %s not found in secret %s\n", key, secretName)
		return "", errors.New("key not found in secret")
	}

	return value, nil
}
