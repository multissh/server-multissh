#!/bin/bash
export CGO_ENABLED=0
echo "Building ssh-client.exe"
GOOS=windows GOARCH=amd64 go build -o ssh-client.exe
echo "Building ssh-client-arm64"
GOOS=linux GOARCH=arm64 go build -o ssh-client-arm64
echo "Building ssh-client-amd64"
GOOS=linux GOARCH=amd64 go build -o ssh-client-amd64
echo "Done!"
