#!/usr/bin/env bash

set -e -o pipefail

if [ -z $VERSION ]; then
    echo "VERSION must be set"
    exit 1
fi

pushd "./$VERSION"

gpg --verify skycoin-${VERSION}-bin-linux-arm.tar.gz.asc
gpg --verify skycoin-${VERSION}-bin-linux-x64.tar.gz.asc
gpg --verify skycoin-${VERSION}-gui-linux-x64.AppImage.asc
gpg --verify skycoin-${VERSION}-bin-win-x64.zip.asc
gpg --verify skycoin-${VERSION}-bin-win-x86.zip.asc
gpg --verify skycoin-${VERSION}-gui-win-setup.exe.asc
gpg --verify skycoin-${VERSION}-bin-osx-darwin-x64.zip.asc
gpg --verify skycoin-${VERSION}-gui-osx-x64.zip.asc
gpg --verify skycoin-${VERSION}-gui-osx.dmg.asc