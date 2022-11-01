#!/bin/bash 

set -e

if [ "$(uname)" == "Darwin" ]; then
    brew install golangci-lint
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1
fi
