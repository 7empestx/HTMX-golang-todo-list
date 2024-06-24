#!/bin/bash
set -e

echo "Go version:"
go version

echo "Current directory structure:"
ls -R

echo "Building Go application..."
TEMPL_EXPERIMENT=rawgo templ generate && go build -o bin/application ./cmd/server/main.go
