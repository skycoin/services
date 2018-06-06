#!/usr/bin/env bash

curl -o gz-c.asc https://raw.githubusercontent.com/skycoin/skycoin/develop/gz-c.asc
echo 'Install the key'
gpg --import gz-c.asc

KEY_ID=ed25519

echo 'Set trust level of the installed key'

echo "$( \
  gpg --list-keys --fingerprint \
  | grep $KEY_ID -A 1 | tail -1 \
  | tr -d '[:space:]' \
):5:" | gpg --import-ownertrust;