#!/bin/bash

set -e

PROJECT_DIR="$(go env GOPATH)/src/github.com/datianshi/simple-cf-gtm"
BUILD_DIR="$PROJECT_DIR/build"

mkdir -p $BUILD_DIR

pushd $PROJECT_DIR
  go test ./...
popd

go test ./...
pushd $PROJECT_DIR/command/
  GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/gtm-linux
  GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/gtm-mac
  GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/gtm-win64
popd
