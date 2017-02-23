#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if [ -d _build ]; then
  rm -rv _build
fi

mkdir _build

OS="linux darwin windows"
ARCH="amd64"

for os in $OS; do
  ext=
  if [ "$os" == "windows" ]; then
    ext=".exe"
  fi

  for arch in $ARCH; do
    echo "Building os=$os arch=$arch"
    file="_build/linky-$os-$arch$ext"
    GOOS=$os GOARCH=$arch go build -o "$file" .
  done
done
