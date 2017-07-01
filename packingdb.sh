#!/bin/sh

set -e
export GOPATH="$(pwd)"
go install github.com/ywwg/packingdb
./bin/packingdb "$@" |less

