#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status.

echo "Starting build and deployment script..."

# Verify that AWS CLI is installed
if ! command -v aws &> /dev/null
then
    echo "AWS CLI could not be found. Please install it to proceed."
    exit 1
fi

# Verify that CDK CLI is installed
if ! command -v cdk &> /dev/null
then
    echo "AWS CDK CLI could not be found. Please install it to proceed."
    exit 1
fi

echo "Building the Go application..."
cd GoApp 

# Check if the main Go file exists
if [ ! -f cmd/server/main.go ]; then
    echo "Error: main.go file not found in cmd/server/"
    exit 1
fi

# Build the Go application
go build -o bin/application cmd/server/main.go
echo "Go application built successfully."
cd ..

# Activate the EB CLI virtual environment
if [ -f ~/.ebcli-virtual-env/bin/activate ]; then
    source ~/.ebcli-virtual-env/bin/activate
else
    echo "EB CLI virtual environment not found. Please set it up to proceed."
    exit 1
fi

# Deploy using EB CLI
echo "Deploying the application using EB CLI..."
eb deploy
echo "Application deployed successfully using EB CLI."

echo "Synthesizing the CDK application..."
cd cdk

# Clean npm cache and install dependencies
echo "Cleaning npm cache and installing dependencies..."
rm -rf node_modules package-lock.json
npm ci

# Synthesize the CDK application
cdk synth
echo "CDK application synthesized successfully."

echo "Deploying the CDK stack..."
cdk deploy --require-approval never
echo "CDK stack deployed successfully."

echo "Invalidating the CloudFront cache..."
aws cloudfront create-invalidation --distribution-id E3UZMP1RS11EPV --paths "/*" | cat

echo "Build and deployment script completed."
