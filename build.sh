#!/bin/bash

set -e

VERSION=$(git describe --tags --always --dirty)
LDFLAGS="-s -w -X main.version=$VERSION"

GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o bin/solvm-linux-amd64
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o bin/solvm-linux-arm64
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o bin/solvm-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o bin/solvm-darwin-arm64
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o bin/solvm-windows-amd64.exe