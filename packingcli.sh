#!/bin/sh

set -e
export GOPATH="$(pwd)"
go install github.com/ywwg/packingcli
./bin/packingcli "$@" | less -X -F
#./bin/packingcli "$@"

