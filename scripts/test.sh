#!/bin/bash -e

set -euf -o pipefail

user_id=$(id -u $(whoami))

docker run --rm \
    -v "$PWD":/go/src/github.com/ankyra/escape-client \
    -w /go/src/github.com/ankyra/escape-client \
    golang:1.8 bash -c "(useradd --uid $user_id builder || true) && su builder -p -c \"set -euf ; /usr/local/go/bin/go test -cover -v \\\$(/usr/local/go/bin/go list ./... | grep -v -E 'vendor' )\""

