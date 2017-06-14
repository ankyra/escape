#!/bin/bash -e

set -e -o pipefail


PLATFORMS="darwin linux"
ARCHS="386 amd64"

echo "$INPUT_credentials" > service_account.json
gcloud auth activate-service-account --key-file service_account.json

echo "Escape v$INPUT_version\n<ul>" > index.html

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
        gcs_target="gs://$INPUT_bucket/escape-client/$INPUT_version/$target"
        gsutil cp "$target" "$gcs_target"
        gsutil acl ch -u AllUsers:R "$gcs_target"
        public_url="https://storage.googleapis.com/$INPUT_bucket/escape-client/$INPUT_version/$target"
        echo "<li><a href=\"$public_url\">$target</a></li>" >> index.html
    done
done

echo "</ul>" >> index.html

gcs_index="gs://$INPUT_bucket/escape-client/$INPUT_version/index.html"
gsutil cp "index.html" "$gcs_index"
gsutil acl ch -u AllUsers:R "$gcs_index"
echo "Published $gcs_index"

