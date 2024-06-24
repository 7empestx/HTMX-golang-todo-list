#!/bin/bash
set -e

echo "Go version:"
go version

echo "Current directory structure:"
ls -R

echo "Building Go application..."
go build -o bin/application ./cmd/server/main.go
