#!/bin/bash -e

set -e -o pipefail

PLATFORMS="darwin linux"
ARCHS="386 amd64"

echo "$INPUT_credentials" > service_account.json

if [ "$INPUT_do_upload" = "1" ] ; then 
    gcloud auth activate-service-account --key-file service_account.json
fi

for GOOS in $PLATFORMS; do
    for ARCH in $ARCHS; do
        target="escape-v$INPUT_escape_version-$GOOS-$ARCH.tgz"
        if [ ! -f $target ] ; then
            docker run --rm -v "$PWD":/go/src/github.com/ankyra/escape \
                            -w /go/src/github.com/ankyra/escape \
                            -e GOOS=$GOOS \
                            -e GOARCH=$ARCH \
                            golang:1.8 go build -v -o escape-$GOOS-$ARCH
            mv escape-${GOOS}-${ARCH} escape
            tar -cvzf ${target} escape
            rm escape
        else
            echo "File $target already exists"
        fi

        if [ "$INPUT_do_upload" = "1" ] ; then 
            gcs_target="gs://$INPUT_bucket/escape/$INPUT_escape_version/$target"
            gsutil cp "$target" "$gcs_target"
            gsutil acl ch -u AllUsers:R "$gcs_target"
            public_url="https://storage.googleapis.com/$INPUT_bucket/escape-client/$INPUT_escape_version/$target"
        fi
    done
done
