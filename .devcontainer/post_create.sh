#!/bin/bash

brew install golang-migrate mockery opentofu gh httpie pgcli codex
# golangci-lint **should not** be installed via brew since it can cause brew to install a different version
# of Go
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.9.0
