#!/usr/bin/env bash

set -e -o pipefail

BASE_URL="https://downloads.skycoin.net/wallet/"

if [ -z $VERSION ]; then
    echo "VERSION must be set"
    exit 1
fi

mkdir -p "$VERSION"
pushd "$VERSION"

curl -o skycoin-${VERSION}-bin-linux-arm.tar.gz ${BASE_URL}skycoin-${VERSION}-bin-linux-arm.tar.gz
curl -o skycoin-${VERSION}-bin-linux-x64.tar.gz ${BASE_URL}skycoin-${VERSION}-bin-linux-x64.tar.gz
curl -o skycoin-${VERSION}-gui-linux-x64.AppImage ${BASE_URL}skycoin-${VERSION}-gui-linux-x64.AppImage
curl -o skycoin-${VERSION}-bin-win-x64.zip ${BASE_URL}skycoin-${VERSION}-bin-win-x64.zip
curl -o skycoin-${VERSION}-bin-win-x86.zip ${BASE_URL}skycoin-${VERSION}-bin-win-x86.zip
curl -o skycoin-${VERSION}-gui-win-setup.exe ${BASE_URL}skycoin-${VERSION}-gui-win-setup.exe
curl -o skycoin-${VERSION}-bin-osx-darwin-x64.zip ${BASE_URL}skycoin-${VERSION}-bin-osx-darwin-x64.zip
curl -o skycoin-${VERSION}-gui-osx-x64.zip ${BASE_URL}skycoin-${VERSION}-gui-osx-x64.zip
curl -o skycoin-${VERSION}-gui-osx.dmg ${BASE_URL}skycoin-${VERSION}-gui-osx.dmg
