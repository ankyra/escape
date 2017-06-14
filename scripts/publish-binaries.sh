#!/bin/bash -e

set -e -o pipefail


PLATFORMS="darwin linux"
ARCHS="386 amd64"

echo "$INPUT_credentials" > service_account.json
gcloud auth activate-service-account --key-file service_account.json

echo "Escape v$INPUT_version\n<ul>" > downloads.html

for GOOS in $PLATFORMS; do
    for ARCH in $ARCHS; do
        target="escape-v$INPUT_version-$GOOS-$ARCH.tgz"
        if [ ! -f $target ] ; then
            docker run --rm -v "$PWD":/go/src/github.com/ankyra/escape-client \
                            -w /go/src/github.com/ankyra/escape-client \
                            -e GOOS=$GOOS \
                            -e GOARCH=$ARCH \
                            golang:1.8 go build -v -o escape-$GOOS-$ARCH
            mv escape-${GOOS}-${ARCH} escape
            tar -cvzf ${target} escape
            rm escape
        else
            echo "File $target already exists"
        fi
        gsutil cp "$target" "gs://$INPUT_bucket/$target"
        gsutil acl ch -u AllUsers:R "gs://$INPUT_bucket/$target"
        public_url="https://storage.googleapis.com/$INPUT_bucket/$target"
        echo "<li><a href=\"$public_url\">$target</a></li>" > downloads.html
    done
done

echo "</ul>" > downloads.html
gsutil cp "downloads.html" "gs://$INPUT_bucket/downloads.html"
gsutil acl ch -u AllUsers:R "gs://$INPUT_bucket/downloads.html"
echo "Published https://storage.googleapis.com/$INPUT_bucket/$target"
