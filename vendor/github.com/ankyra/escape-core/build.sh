#!/bin/bash -e

set -euf -o pipefail

user_id=$(id -u $(whoami))

docker run --rm \
    -v "$PWD":/go/src/github.com/ankyra/escape-core \
    -w /go/src/github.com/ankyra/escape-core \
    golang:1.8 bash -c "(useradd --uid $user_id builder || true) && su builder -p -c \"/usr/local/go/bin/go build\""
