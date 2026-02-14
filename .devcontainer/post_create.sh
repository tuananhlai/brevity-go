#!/bin/bash

brew install golang-migrate mockery
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.9.0