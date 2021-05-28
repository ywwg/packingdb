#!/bin/sh

set -e
export GOPATH="$(pwd)"
export GO111MODULE="off"
go install github.com/ywwg/packingdb
./bin/packingdb "$@"

