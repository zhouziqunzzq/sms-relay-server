#!/bin/bash
set -e

# Set Go environment variables for cross-compiling to Linux/amd64
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

go build -tags lambda.norpc -o bootstrap

# Package the binary for AWS Lambda (assumes build-lambda-zip is in PATH)
build-lambda-zip -o sms-relay-api-authenticator.zip bootstrap

