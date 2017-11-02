#!/bin/bash -e

set -e -o pipefail


PLATFORMS="darwin linux"
ARCHS="386 amd64"

echo "$INPUT_credentials" > service_account.json
gcloud auth activate-service-account --key-file service_account.json

echo "<h2>Escape v$INPUT_escape_version</h2><ul>" > index.html
echo "<h2>Escape v$INPUT_escape_version</h2><ul>" > main_index.html

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
        gcs_target="gs://$INPUT_bucket/escape/$INPUT_escape_version/$target"
        gsutil cp "$target" "$gcs_target"
        gsutil acl ch -u AllUsers:R "$gcs_target"
        public_url="https://storage.googleapis.com/$INPUT_bucket/escape-client/$INPUT_escape_version/$target"
        echo "<li><a href=\"$public_url\">$target</a></li>" >> index.html
        echo "<li><a href=\"$public_url\">$target</a></li>" >> main_index.html
    done
done

echo "</ul>" >> index.html

gcs_index="gs://$INPUT_bucket/escape-client/$INPUT_escape_version/index.html"
gsutil cp "index.html" "$gcs_index"
gsutil acl ch -u AllUsers:R "$gcs_index"
echo "Published $gcs_index"

echo "</ul><h3>Older releases</h3><ul class='older-releases'>" >> main_index.html
gsutil ls 'gs://escape-releases-eu/escape-client/' | while read line ; do \
    echo $line | sed 's|gs://escape-releases-eu/escape-client/\(.*\)/|<li><a href="https://storage.googleapis.com/'$INPUT_bucket'/escape-client/\1/index.html">\1</a></li>|' | tee --append main_index.html
done
echo "</ul>" >> main_index.html

gcs_index="gs://$INPUT_bucket/escape-client/index.html"
gsutil cp "main_index.html" "$gcs_index"
gsutil acl ch -u AllUsers:R "$gcs_index"
echo "Updated $gcs_index"
