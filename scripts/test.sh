#!/bin/bash -e

set -euf -o pipefail

user_id=$(id -u $(whoami))

docker rm src || true
docker create -v /go/src/github.com/ankyra/ --name src golang:1.9.0 /bin/true
docker cp "$PWD" src:/go/src/github.com/ankyra/tmp
docker run --rm --volumes-from src \
    -w /go/src/github.com/ankyra/ \
    golang:1.9.0 mv tmp escape-client
docker run --rm \
    --volumes-from src \
    -w /go/src/github.com/ankyra/escape-client \
    golang:1.9.0 bash -c "go test -cover -v \$(go list ./... | grep -v -E 'vendor|godog' )"
docker rm src
