#!/bin/bash
set -eu
[ -f "$(go env GOPATH)/bin/gotests" ] || go install -v github.com/cweill/gotests/gotests@latest
[ -f "$(go env GOPATH)/bin/gomodifytags" ] || go install -v github.com/fatih/gomodifytags@latest
[ -f "$(go env GOPATH)/bin/impl" ] || go install -v github.com/josharian/impl@latest
[ -f "$(go env GOPATH)/bin/goplay" ] || go install -v github.com/haya14busa/goplay/cmd/goplay@latest
[ -f "$(go env GOPATH)/bin/dlv" ] || go install -v github.com/go-delve/delve/cmd/dlv@latest
[ -f "$(go env GOPATH)/bin/golangci-lint" ] || go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
[ -f "$(go env GOPATH)/bin/gopls" ] || go install -v golang.org/x/tools/gopls@latest
