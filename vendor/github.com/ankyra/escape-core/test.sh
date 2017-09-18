#!/bin/bash -e

set -euf -o pipefail

user_id=$(id -u $(whoami))
username="builder"
if test $user_id '==' 0 ; then 
  username="root"
fi

docker run --rm \
    -v "$PWD":/go/src/github.com/ankyra/escape-core \
    -w /go/src/github.com/ankyra/escape-core \
    golang:1.9.0 bash -c "(test $user_id '==' 0 || useradd --uid $user_id builder || true) && su $username -p -c \"set -euf ; /usr/local/go/bin/go test -cover -v \\\$(/usr/local/go/bin/go list ./... | grep -v -E 'vendor|godog' )\""

