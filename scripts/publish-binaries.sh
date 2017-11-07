#!/bin/bash -e

set -e -o pipefail

PLATFORMS="darwin linux"
ARCHS="386 amd64"

BASE_DIR=$(dirname "$(readlink -f "$0")")
SRC_DIR=$(readlink -f "${BASE_DIR}/../")
GOLANG_VERSION=1.9.0
BUILD_IMAGE="golang:${GOLANG_VERSION}"

echo "$INPUT_credentials" > service_account.json

if [ "$INPUT_do_upload" = "1" ] ; then 
    gcloud auth activate-service-account --key-file service_account.json
fi

for GOOS in $PLATFORMS; do
    for ARCH in $ARCHS; do
        filename="escape-v$INPUT_escape_version-$GOOS-$ARCH.tgz"
        target="${SRC_DIR}/${filename}"
        if [ ! -f $target ] ; then
            echo "Building for $GOOS-$ARCH from ${SRC_DIR}"
            docker rm src || true
            docker create -v /go/src/github.com/ankyra/ --name src ${BUILD_IMAGE} /bin/true
            docker cp "$SRC_DIR" src:/go/src/github.com/ankyra/tmp
            docker run --rm --volumes-from src \
                -w /go/src/github.com/ankyra/ \
                ${BUILD_IMAGE} mv tmp escape
            docker run --rm \
                --volumes-from src \
                -w /go/src/github.com/ankyra/escape \
                -e GOOS=$GOOS \
                -e GOARCH=$ARCH \
                ${BUILD_IMAGE} bash -c "go build -v -o escape-$GOOS-$ARCH"
            docker cp src:/go/src/github.com/ankyra/escape/escape-${GOOS}-${ARCH} ${SRC_DIR}/escape
            docker rm src
            echo "Creating archive: ${target}"
            tar -C "${SRC_DIR}" -cvzf "${target}" "escape"
            rm -f "${SRC_DIR}/escape"
        else
            echo "File $target already exists"
        fi

        if [ "$INPUT_do_upload" = "1" ] ; then 
            gcs_target="gs://$INPUT_bucket/escape/$INPUT_escape_version/$filename"
            echo "Copying to $gcs_target"
            gsutil cp "$target" "$gcs_target"
            echo "Setting ACL on $gcs_target"
            gsutil acl ch -u AllUsers:R "$gcs_target"
        fi
    done
done

rm service_account.json
