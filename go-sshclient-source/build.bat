@echo off
SET CGO_ENABLED=0
setlocal enabledelayedexpansion
SET GOOS=windows
SET GOARCH=amd64
echo "Building ssh-client.exe"
go build -o ssh-client.exe
SET GOOS=linux
SET GOARCH=arm64
echo "Building ssh-client-arm64"
go build -o ssh-client-arm64
SET GOOS=linux
SET GOARCH=amd64
echo "Building ssh-client-amd64"
go build -o ssh-client-amd64
echo "Done!"
