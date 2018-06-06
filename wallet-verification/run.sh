#!/usr/bin/env bash

sh ./download-builds.sh
sh ./download-signatures.sh
sh ./install-public-key.sh
sh ./verify-signatures.sh

echo "The signatures have been successfully verified."