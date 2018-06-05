#!/usr/bin/env bash

curl -o gz-c.asc https://raw.githubusercontent.com/skycoin/skycoin/develop/gz-c.asc
gpg --import gz-c.asc