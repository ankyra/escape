#!/bin/bash -e

set -euf -o pipefail

cat > util/metadata.go <<EOF
package util

const EscapeVersion="$INPUT_version"
EOF
