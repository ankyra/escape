#!/bin/bash -e

set -euf -o pipefail

user_id=$(id -u $(whoami))

rm -rf vendor/github.com/ankyra/escape-core
cp -r deps/_/escape-core/ vendor/github.com/ankyra/escape-core
rm -rf vendor/github.com/ankyra/escape-core/vendor/

docker run --rm \
    -v "$PWD":/go/src/github.com/ankyra/escape-client \
    -w /go/src/github.com/ankyra/escape-client \
    golang:1.8 bash -c "(useradd --uid $user_id builder || true) && su builder -p -c \"/usr/local/go/bin/go build -v -o escape && mkdir -p docs/cmd && /usr/local/go/bin/go run docs/generate_cmd_docs.go\""
