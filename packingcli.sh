#!/bin/sh

set -e
export GOPATH="$(pwd)"
go install github.com/ywwg/cmd/packingcli
#./bin/packingcli "$@" | less -X -F
./bin/packingcli "$@"

