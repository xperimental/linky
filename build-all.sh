#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if [ -d _build ]; then
  rm -rv _build
fi

mkdir _build

readonly OS="linux darwin windows"
readonly ARCH="amd64"
VERSION="$(git describe --tags)"; readonly VERSION

for os in $OS; do
  ext=
  if [ "$os" == "windows" ]; then
    ext=".exe"
  fi

  for arch in $ARCH; do
    echo "Building os=$os arch=$arch"
    file="_build/linky-$VERSION-$os-$arch$ext"
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o "$file" -ldflags "-w" .
    upx -9 "$file"
  done
done
