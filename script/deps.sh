#!/bin/bash

NEEDED_COMMANDS="go git goreleaser golangci-lint"

for cmd in ${NEEDED_COMMANDS} ; do
    if ! command -v "${cmd}" &> /dev/null ; then
        echo -e "\033[91m${cmd} missing please install \033[0m"
        exit 1
    else
        echo "${cmd} found"
    fi
done

# upx
## https://upx.github.io/

# Git
## https://git-scm.com/

# Go
## https://golang.org/

# goreleaser
## https://github.com/goreleaser/goreleaser
## go install github.com/goreleaser/goreleaser
## Mac: brew install goreleaser
## Ubuntu: apt install goreleaser

# golang-ci-lint
## https://github.com/golangci/golangci-lint
## Mac: brew install golangci-lint
## go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1