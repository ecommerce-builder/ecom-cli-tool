#!/bin/bash
echo $ECOM_CLI_VERSION
GOOS=linux   GOARCH=amd64 go build -ldflags "-X main.version=$ECOM_CLI_VERSION" -o build/ecom-$ECOM_CLI_VERSION
GOOS=darwin  GOARCH=amd64 go build -ldflags "-X main.version=$ECOM_CLI_VERSION" -o build/ecom-$ECOM_CLI_VERSION-darwin
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$ECOM_CLI_VERSION" -o build/ecom-$ECOM_CLI_VERSION-windows
