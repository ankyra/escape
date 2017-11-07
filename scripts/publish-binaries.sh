#!/bin/bash -e

set -e -o pipefail

PLATFORMS="darwin linux"
ARCHS="386 amd64"

BASE_DIR=$(dirname "$(readlink -f "$0")")
SRC_DIR=$(readlink -f "${BASE_DIR}/../")

echo "$INPUT_credentials" > service_account.json

if [ "$INPUT_do_upload" = "1" ] ; then 
    gcloud auth activate-service-account --key-file service_account.json
fi

for GOOS in $PLATFORMS; do
    for ARCH in $ARCHS; do
        filename="escape-v$INPUT_escape_version-$GOOS-$ARCH.tgz"
        target="${SRC_DIR}/${filename}"
        if [ ! -f $target ] ; then
            docker run --rm -v "$SRC_DIR":/go/src/github.com/ankyra/escape \
                            -w /go/src/github.com/ankyra/escape \
                            -e GOOS=$GOOS \
                            -e GOARCH=$ARCH \
                            golang:1.8 go build -v -o escape-$GOOS-$ARCH
            mv escape-${GOOS}-${ARCH} "${SRC_DIR}/escape"
            tar -C "${SRC_DIR}" -cvzf "${target}" "escape"
            rm escape
        else
            echo "File $target already exists"
        fi

        if [ "$INPUT_do_upload" = "1" ] ; then 
            gcs_target="gs://$INPUT_bucket/escape/$INPUT_escape_version/$filename"
            gsutil cp "$target" "$gcs_target"
            gsutil acl ch -u AllUsers:R "$gcs_target"
        fi
    done
done
