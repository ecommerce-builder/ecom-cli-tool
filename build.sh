#!/bin/bash
echo $ECOM_CLI_VERSION
echo $ECOM_DOCS_PUBLIC_DOWNLOAD
GOOS=linux   GOARCH=amd64 go build -ldflags "-X main.version=$ECOM_CLI_VERSION" -o build/ecom-$ECOM_CLI_VERSION-linux
GOOS=darwin  GOARCH=amd64 go build -ldflags "-X main.version=$ECOM_CLI_VERSION" -o build/ecom-$ECOM_CLI_VERSION-darwin
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$ECOM_CLI_VERSION" -o build/ecom-$ECOM_CLI_VERSION.exe

# Linux
rm -f $ECOM_DOCS_PUBLIC_DOWNLOAD/linux/ecom*
cp build/ecom-$ECOM_CLI_VERSION-linux $ECOM_DOCS_PUBLIC_DOWNLOAD/linux/ecom-$ECOM_CLI_VERSION-linux


# Mac
rm -f $ECOM_DOCS_PUBLIC_DOWNLOAD/mac/ecom*
cp build/ecom-$ECOM_CLI_VERSION-darwin $ECOM_DOCS_PUBLIC_DOWNLOAD/mac/ecom-$ECOM_CLI_VERSION-darwin

# Windows
rm -f $ECOM_DOCS_PUBLIC_DOWNLOAD/windows/ecom*
cp build/ecom-$ECOM_CLI_VERSION.exe $ECOM_DOCS_PUBLIC_DOWNLOAD/windows/ecom-$ECOM_CLI_VERSION.exe
