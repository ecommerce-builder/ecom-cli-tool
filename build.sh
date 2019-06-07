#!/bin/bash
echo $VERSION
GOOS=linux   GOARCH=amd64 go build -o build/ecom-$VERSION-linux
GOOS=darwin  GOARCH=amd64 go build -o build/ecom-$VERSION-darwin
GOOS=windows GOARCH=amd64 go build -o build/ecom-$VERSION.exe
